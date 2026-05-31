package kun

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"motewallet/internal/config"
	kundto "motewallet/internal/pkg/kun/dto"
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
	return m.PostAsCustomer(ctx, m.customerNo, path, reqBody, respBody)
}

func (m *MockClient) PostAsCustomer(ctx context.Context, customerNo, path string, reqBody interface{}, respBody interface{}) error {
	reqJSON, _ := json.Marshal(reqBody)
	log.Printf("[KUN Mock] POST %s Customer-No=%s body=%s", path, customerNo, string(reqJSON))

	var mockData []byte

	switch path {
	case "/rest/v2.0/customer/register":
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		mockData = []byte(fmt.Sprintf(`{"subCustomerNo":"MOCK_SUB_%s"}`, hex.EncodeToString(b)))
	case "/rest/v2.0/customer/agreeList":
		mockData = []byte(`{"list":[]}`)
	case "/rest/v2.0/customer/subMerchant/register":
		mockData = []byte(`{"authId":"MOCK_AUTH_001"}`)
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

func (m *MockClient) UploadFileAsCustomer(
	ctx context.Context,
	customerNo, filename string,
	content []byte,
	contentType string,
) (*kundto.FileUploadResp, error) {
	log.Printf("[KUN Mock] POST /rest/v2.0/upload Customer-No=%s filename=%s size=%d", customerNo, filename, len(content))
	return &kundto.FileUploadResp{
		URL: fmt.Sprintf("https://mock.kun.global/files/%s", hex.EncodeToString([]byte(filename))[:12]),
	}, nil
}
