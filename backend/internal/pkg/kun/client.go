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

	"motewallet-withdrawal/backend/internal/config"
)

// KUNClient defines the interface for KUN API calls.
type KUNClient interface {
	Post(ctx context.Context, path string, reqBody interface{}, respBody interface{}) error
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
	apiKey     string
	apiSecret  string
	customerNo string
	regionCode string
	httpClient *http.Client
}

func NewClient(cfg config.KUNConfig) *Client {
	return &Client{
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		apiSecret:  cfg.APISecret,
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
	sign := SignRequest(params, timestamp, c.apiSecret)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", c.apiKey)
	req.Header.Set("sign", sign)
	req.Header.Set("timestamp", timestamp)

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
