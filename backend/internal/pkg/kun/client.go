package kun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"motewallet/internal/config"
	kundto "motewallet/internal/pkg/kun/dto"
)

// KUNClient defines the interface for KUN API calls.
type KUNClient interface {
	Post(ctx context.Context, path string, reqBody interface{}, respBody interface{}) error
	// PostAsCustomer signs the request with the given Customer-No (e.g. sub-merchant no for onboarding auth).
	PostAsCustomer(ctx context.Context, customerNo, path string, reqBody interface{}, respBody interface{}) error
	// UploadFileAsCustomer uploads a file via POST /rest/v2.0/upload (multipart/form-data).
	// See: https://opendocs.kun.global/docs/api/file-upload
	UploadFileAsCustomer(ctx context.Context, customerNo, filename string, content []byte, contentType string) (*kundto.FileUploadResp, error)
	GetRegionCode() string
	GetCustomerNo() string
}

type kunResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// Client implements KUNClient with real HTTP calls.
type Client struct {
	baseURL    string
	appKey     string
	apiVersion string
	privateKey string
	publicKey  string
	customerNo string
	regionCode string
	httpClient *http.Client
}

func NewClient(cfg config.KUNConfig) *Client {
	return &Client{
		baseURL:    cfg.BaseURL,
		appKey:     cfg.AppKey,
		apiVersion: cfg.ApiVersion,
		privateKey: cfg.PrivateKey,
		publicKey:  cfg.PublicKey,
		customerNo: cfg.CustomerNo,
		regionCode: cfg.RegionCode,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *Client) GetRegionCode() string {
	return c.regionCode
}

func (c *Client) GetCustomerNo() string {
	return c.customerNo
}

func (c *Client) Post(ctx context.Context, path string, reqBody interface{}, respBody interface{}) error {
	return c.PostAsCustomer(ctx, c.customerNo, path, reqBody, respBody)
}

func (c *Client) PostAsCustomer(ctx context.Context, customerNo, path string, reqBody interface{}, respBody interface{}) error {
	bodyBytes, err := MarshalRequestBody(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	params, err := mapFromJSON(bodyBytes)
	if err != nil {
		return fmt.Errorf("flatten request body: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signString, sign, err := CanonicalizeAndSignRequest(params, customerNo, timestamp, c.privateKey)
	if err != nil {
		return fmt.Errorf("sign request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Key", c.appKey)
	req.Header.Set("Api-Version", c.apiVersion)
	req.Header.Set("Customer-No", customerNo)
	req.Header.Set("Timestamp", timestamp)
	req.Header.Set("Sign", sign)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	logKUNRequest(http.MethodPost, url, path, signString, req.Header, resp.Header, json.RawMessage(bodyBytes), resp.StatusCode, nil)

	var kunResp kunResponse
	if err := json.Unmarshal(respBytes, &kunResp); err != nil {
		return fmt.Errorf("unmarshal KUN response: %w", err)
	}

	// Verify response signature when provided.
	// Doc: server may return sign+timestamp in headers (or body); we validate when headers exist.
	if c.publicKey != "" {
		respSign := resp.Header.Get("Sign")
		if respSign == "" {
			respSign = resp.Header.Get("sign")
		}
		respTimestamp := resp.Header.Get("Timestamp")
		if respTimestamp == "" {
			respTimestamp = resp.Header.Get("timestamp")
		}
		if respSign != "" && respTimestamp != "" && len(kunResp.Data) > 0 {
			var dataAny interface{}
			_ = json.Unmarshal(kunResp.Data, &dataAny)
			dataMap := map[string]interface{}{}
			switch t := dataAny.(type) {
			case map[string]interface{}:
				dataMap = t
			default:
				// If data is not an object, verify under a stable key.
				dataMap["data"] = t
			}
			if err := VerifyResponseSignature(c.publicKey, dataMap, respTimestamp, respSign); err != nil {
				return fmt.Errorf("verify KUN response signature: %w", err)
			}
		}
	}

	if kunResp.Code != "200" && kunResp.Code != "000000" {
		return &KUNError{
			Code:    kunResp.Code,
			Message: kunResp.Msg,
		}
	}

	if respBody != nil && len(kunResp.Data) > 0 {
		if err := json.Unmarshal(kunResp.Data, respBody); err != nil {
			return fmt.Errorf("unmarshal response data: %w", err)
		}
	}

	return nil
}

func (c *Client) UploadFileAsCustomer(
	ctx context.Context,
	customerNo, filename string,
	content []byte,
	contentType string,
) (*kundto.FileUploadResp, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(content); err != nil {
		return nil, fmt.Errorf("write file content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	url := c.baseURL + "/rest/v2.0/upload"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signString, sign, err := CanonicalizeAndSignRequest(map[string]interface{}{}, customerNo, timestamp, c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("App-Key", c.appKey)
	req.Header.Set("Api-Version", c.apiVersion)
	req.Header.Set("Customer-No", customerNo)
	req.Header.Set("Timestamp", timestamp)
	req.Header.Set("Sign", sign)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	uploadPath := "/rest/v2.0/upload"
	logKUNRequest(http.MethodPost, url, uploadPath, signString, req.Header, resp.Header, map[string]any{
		"type":         "multipart",
		"filename":     filename,
		"content_type": contentType,
		"size":         len(content),
	}, resp.StatusCode, respBytes)

	var kunResp struct {
		Code string                `json:"code"`
		Msg  string                `json:"message"`
		Data kundto.FileUploadResp `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &kunResp); err != nil {
		return nil, fmt.Errorf("unmarshal KUN response: %w", err)
	}

	if kunResp.Code != "000000" && kunResp.Code != "200" {
		return nil, &KUNError{
			Code:    kunResp.Code,
			Message: kunResp.Msg,
		}
	}

	if kunResp.Data.URL == "" {
		return nil, fmt.Errorf("KUN file upload returned empty url")
	}

	return &kunResp.Data, nil
}
