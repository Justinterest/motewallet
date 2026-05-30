package kun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"motewallet/internal/config"
)

// KUNClient defines the interface for KUN API calls.
type KUNClient interface {
	Post(ctx context.Context, path string, reqBody interface{}, respBody interface{}) error
	// PostAsCustomer signs the request with the given Customer-No (e.g. sub-merchant no for onboarding auth).
	PostAsCustomer(ctx context.Context, customerNo, path string, reqBody interface{}, respBody interface{}) error
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
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	params, err := StructToMap(reqBody)
	if err != nil {
		return fmt.Errorf("flatten request body: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	_, sign, err := CanonicalizeAndSignRequest(params, customerNo, timestamp, c.privateKey)
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

	if kunResp.Code != "200" {
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
