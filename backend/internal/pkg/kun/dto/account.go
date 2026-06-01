package dto

// DepositAddressListReq is the body for POST /rest/v2.0/customer/crypto/deposit/addresses.
// See: https://opendocs.kun.global/docs/api/crypto-deposit-address-list-query
type DepositAddressListReq struct {
	RequestNo string `json:"requestNo"`
	Currency  string `json:"currency"`
	ChainType string `json:"chainType"`
}

type DepositAddressItem struct {
	Address   string `json:"address"`
	Currency  string `json:"currency"`
	ChainType string `json:"chainType"`
	Chain     string `json:"chain"`
}

// DepositHistoryQueryReq is the body for POST /rest/v2.0/trade/digital/wallet/query/recharge.
// See: https://opendocs.kun.global/docs/api/crypto-deposit-history-query
type DepositHistoryQueryReq struct {
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	OrderCurrency string `json:"orderCurrency,omitempty"`
	Chain         string `json:"chain,omitempty"`
	OrderID       string `json:"orderId,omitempty"`
	PageNo        int    `json:"pageNo"`
	PageSize      int    `json:"pageSize"`
}

type DepositHistoryRecord struct {
	OrderID              string `json:"orderId"`
	CustomerID           string `json:"customerId"`
	OrderCurrency        string `json:"orderCurrency"`
	OrderAmount          string `json:"orderAmount"`
	Chain                string `json:"chain"`
	TxID                 string `json:"txId"`
	NetworkConfirmNumber string `json:"networkConfirmNumber"`
	OrderTime            string `json:"orderTime"`
	UpdateTime           string `json:"updateTime"`
	FeeCurrency          string `json:"feeCurrency"`
	FeeAmount            string `json:"feeAmount"`
	OrderStatus          string `json:"orderStatus"`
	SendWalletAddress    string `json:"sendWalletAddress"`
	ReceiveWalletAddress string `json:"receiveWalletAddress"`
	ErrorCode            string `json:"errorCode,omitempty"`
	ErrorMessage         string `json:"errorMessage,omitempty"`
}

type DepositHistoryQueryResp struct {
	PageNo      int                    `json:"pageNo"`
	PageSize    int                    `json:"pageSize"`
	PageCount   int                    `json:"pageCount"`
	RecordCount int                    `json:"recordCount"`
	Records     []DepositHistoryRecord `json:"records"`
	TotalSize   int                    `json:"totalSize"`
	TotalPage   int                    `json:"totalPage"`
	Rows        []DepositHistoryRecord `json:"rows"`
}

func (r *DepositHistoryQueryResp) Items() []DepositHistoryRecord {
	if len(r.Rows) > 0 {
		return r.Rows
	}
	return r.Records
}

func (r *DepositHistoryQueryResp) Total() int64 {
	if r.TotalSize > 0 {
		return int64(r.TotalSize)
	}
	if r.RecordCount > 0 {
		return int64(r.RecordCount)
	}
	return int64(len(r.Items()))
}

type BalanceQueryReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	Currency      string `json:"currency"`
	CurrencyType  string `json:"currencyType,omitempty"`
}

type BalanceQueryResp struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Total     string `json:"total"`
}

type CryptoAddressAddReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	Currency      string `json:"currency"`
	Chain         string `json:"chain"`
	Address       string `json:"address"`
	Alias         string `json:"alias,omitempty"`
	RequestNo     string `json:"requestNo"`
}

type CryptoAddressAddResp struct {
	AccountId string `json:"accountId"`
}

type FiatAddressAddReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	Currency      string `json:"currency"`
	BankName      string `json:"bankName"`
	BankCountry   string `json:"bankCountry"`
	SwiftCode     string `json:"swiftCode"`
	AccountName   string `json:"accountName"`
	AccountNo     string `json:"accountNo"`
	TransferType  string `json:"transferType"`
	RequestNo     string `json:"requestNo"`
}

type FiatAddressAddResp struct {
	AccountId string `json:"accountId"`
}

type RegionBalanceQueryReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	RegionCode    string `json:"regionCode"`
}

type RegionBalanceItem struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Total     string `json:"total"`
}

type RegionBalanceQueryResp struct {
	Balances []RegionBalanceItem `json:"balances"`
}
