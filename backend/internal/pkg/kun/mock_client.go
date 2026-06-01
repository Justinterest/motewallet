package kun

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
	case "/rest/v2.0/customer/fiat/withdrawal/countries":
		mockData = []byte(`[
			{"countryName":"中国香港","countryCode":"HK"},
			{"countryName":"中国大陆","countryCode":"CN"},
			{"countryName":"新加坡","countryCode":"SG"},
			{"countryName":"美国","countryCode":"US"},
			{"countryName":"英国","countryCode":"GB"}
		]`)
	case "/rest/v2.0/customer/country/auth/types":
		countryCode := "HK"
		if reqBody != nil {
			var req kundto.CountryAuthTypesReq
			if err := json.Unmarshal(reqJSON, &req); err == nil && req.CountryCode != "" {
				countryCode = strings.ToUpper(req.CountryCode)
			}
		}
		switch countryCode {
		case "CN":
			mockData = []byte(`[
				{"docName":"身份证","docCode":"MOCK_CN_ID"},
				{"docName":"护照","docCode":"MOCK_CN_PASSPORT"}
			]`)
		default:
			mockData = []byte(`[
				{"docName":"香港永久性居民身份证","docCode":"MOCK_HK_ID"},
				{"docName":"护照","docCode":"MOCK_HK_PASSPORT"}
			]`)
		}
	case "/rest/v2.0/customer/crypto/deposit/addresses":
		mockData = []byte(`[
			{"address":"TMockDepositAddress123456789","currency":"USDT","chainType":"TRX_TRC20","chain":"Tron"},
			{"address":"0xMockDepositAddress1234567890abcdef","currency":"USDT","chainType":"ETH_ERC20","chain":"Ethereum"}
		]`)
	case "/rest/v2.0/trade/digital/wallet/query/recharge":
		mockData = []byte(`{
			"totalSize": 1,
			"totalPage": 1,
			"rows": [{
				"orderId": "MOCK_DEP_001",
				"orderCurrency": "USDT",
				"orderAmount": "100.00000000",
				"chain": "TRX_TRC20",
				"txId": "mock_tx_id",
				"orderTime": "2025-11-18 10:07:41",
				"orderStatus": "SUCCESS"
			}]
		}`)
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
