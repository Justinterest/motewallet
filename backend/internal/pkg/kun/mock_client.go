package kun

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

	var mockData []byte

	switch path {
	case "/rest/v2.0/customer/register":
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		mockData = []byte(fmt.Sprintf(`{"subCustomerNo":"MOCK_SUB_%s"}`, hex.EncodeToString(b)))
	case "/rest/v2.0/customer/agreeList":
		mockData = []byte(`[
			{"protocolId":"MOCK_PROTO_001","title":"KUN SPACE 服务协议","url":"https://example.com/agreement.pdf","signStatus":"UNSIGN","version":"1.0"},
			{"protocolId":"MOCK_PROTO_002","title":"隐私政策","url":"https://example.com/privacy.pdf","signStatus":"UNSIGN","version":"1.0"}
		]`)
	case "/rest/v2.0/customer/agree/auth":
		mockData = []byte(`{}`)
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
	case "/rest/v2.0/account/query/balance":
		mockData = []byte(`[
			{"currency":"USD","balance":"1000.00000000","regionCode":"KUN_PL"},
			{"currency":"HKD","balance":"5000.00000000","regionCode":"KUN_PL"},
			{"currency":"USDT","balance":"820.00000000","regionCode":"KUN_PL"},
			{"currency":"USDC","balance":"300.00000000","regionCode":"KUN_PL"},
			{"currency":"BTC","balance":"0.50000000","regionCode":"KUN_PL"}
		]`)
	case "/rest/v2.0/user/fund/transfer":
		mockData = []byte(`{"orderId":"MOCK_TRANSFER_001","orderStatus":"SUCCESS"}`)
	case "/rest/v2.0/trade/account/outAccount/query":
		currencyCode := "USDT"
		if reqBody != nil {
			var req kundto.OutAccountBalanceQueryReq
			if err := json.Unmarshal(reqJSON, &req); err == nil && req.Currency != "" {
				currencyCode = strings.ToUpper(req.Currency)
			}
		}
		mockData = []byte(fmt.Sprintf(`[{"currency":"%s","balance":"205.00000000"}]`, currencyCode))
	case "/rest/v2.0/customer/fiat/address/add":
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		mockData = []byte(fmt.Sprintf(`{"accountId":"MOCK_BANK_%s"}`, hex.EncodeToString(b)))
	case "/rest/v2.0/customer/fiat/withdrawal/del":
		mockData = []byte(`{}`)
	case "/rest/v2.0/trade/fiat/withdrawal":
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		mockData = []byte(fmt.Sprintf(`{"orderId":"MOCK_FIAT_WD_%s"}`, hex.EncodeToString(b)))
	case "/rest/v2.0/customer/crypto/address/add":
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		mockData = []byte(fmt.Sprintf(`{"accountId":"MOCK_CRYPTO_ADDR_%s"}`, hex.EncodeToString(b)))
	case "/rest/v2.0/customer/crypto/address/del":
		mockData = []byte(`{}`)
	case "/rest/v2.0/trade/crypto/withdrawal":
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		mockData = []byte(fmt.Sprintf(`"MOCK_CRYPTO_WD_%s"`, hex.EncodeToString(b)))
	case "/rest/v2.0/trade/inner/match/create":
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		mockData = []byte(fmt.Sprintf(`{"orderId":"MOCK_1TO1_%s"}`, hex.EncodeToString(b)))
	case "/rest/v2.0/trade/inner/match/query":
		mockData = []byte(`{
			"orderId":"MOCK_1TO1_001",
			"orderStatus":"SUCCESS",
			"fromCurrency":"USDT",
			"toCurrency":"USD",
			"orderAmount":"100.00000000",
			"orderCurrency":"USDT",
			"receiveAmount":"100.00000000",
			"receiveCurrency":"USD",
			"exchangeRate":"1",
			"tradeFee":"0",
			"tradeFeeCurrency":"USDT"
		}`)
	default:
		mockData = []byte(`{}`)
	}

	logKUNMockCall(path, customerNo, json.RawMessage(reqJSON), mockData)

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
	uploadPath := "/rest/v2.0/upload"
	resp := &kundto.FileUploadResp{
		URL: fmt.Sprintf("https://mock.kun.global/files/%s", hex.EncodeToString([]byte(filename))[:12]),
	}
	data, _ := json.Marshal(resp)
	logKUNMockCall(uploadPath, customerNo, map[string]any{
		"type":         "multipart",
		"filename":     filename,
		"content_type": contentType,
		"size":         len(content),
	}, data)
	return resp, nil
}
