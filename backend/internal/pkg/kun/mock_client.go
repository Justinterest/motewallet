package kun

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"motewallet/internal/config"
)

// MockClient implements KUNClient for development/testing.
type MockClient struct {
	regionCode string
	customerNo string
}

func NewMockClient(cfg config.KUNConfig) *MockClient {
	return &MockClient{
		regionCode: cfg.RegionCode,
		customerNo: cfg.CustomerNo,
	}
}

func (m *MockClient) GetRegionCode() string {
	return m.regionCode
}

func (m *MockClient) GetCustomerNo() string {
	return m.customerNo
}

func (m *MockClient) Post(ctx context.Context, path string, reqBody interface{}, respBody interface{}) error {
	reqJSON, _ := json.Marshal(reqBody)
	log.Printf("[KUN Mock] POST %s body=%s", path, string(reqJSON))

	var mockData []byte

	switch path {
	case "/rest/v2.0/customer/register":
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		mockData = []byte(fmt.Sprintf(`{"subCustomerNo":"MOCK_SUB_%s"}`, hex.EncodeToString(b)))
	case "/rest/v2.0/customer/agreeList":
		mockData = []byte(`{"list":[]}`)
	case "/rest/v2.0/customer/subMerchant/register":
		mockData = []byte(`{}`)
	case "/rest/v2.0/customer/merchant/register/query":
		mockData = []byte(`{"authStatus":"AUTH_SUC"}`)
	default:
		mockData = []byte(`{}`)
	}

	if respBody != nil {
		if err := json.Unmarshal(mockData, respBody); err != nil {
			return fmt.Errorf("mock unmarshal: %w", err)
		}
	}

	return nil
}
