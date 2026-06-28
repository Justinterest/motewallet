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
	PageNo      FlexInt                `json:"pageNo"`
	PageSize    FlexInt                `json:"pageSize"`
	PageCount   FlexInt                `json:"pageCount"`
	RecordCount FlexInt                `json:"recordCount"`
	Records     []DepositHistoryRecord `json:"records"`
	TotalSize   FlexInt                `json:"totalSize"`
	TotalPage   FlexInt                `json:"totalPage"`
	Rows        []DepositHistoryRecord `json:"rows"`
}

func (r *DepositHistoryQueryResp) Items() []DepositHistoryRecord {
	if len(r.Rows) > 0 {
		return r.Rows
	}
	return r.Records
}

func (r *DepositHistoryQueryResp) Total() int64 {
	if r.TotalSize.Int() > 0 {
		return int64(r.TotalSize.Int())
	}
	if r.RecordCount.Int() > 0 {
		return int64(r.RecordCount.Int())
	}
	return int64(len(r.Items()))
}

type BalanceQueryReq struct {
	RequestNo     string `json:"requestNo"`
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

// FiatAddressAddReq is the body for POST /rest/v2.0/customer/fiat/address/add.
// See: https://opendocs.kun.global/docs/api/bind-fiat-withdrawal-account
type FiatAddressAddReq struct {
	RequestNo          string   `json:"requestNo"`
	AccountCategory    string   `json:"accountCategory"`
	CurrencyList       []string `json:"currencyList"`
	Area               string   `json:"area"`
	TransferType       string   `json:"transferType"`
	AccountNo          string   `json:"accountNo"`
	AccountName        string   `json:"accountName"`
	AccountType        string   `json:"accountType,omitempty"`
	SwiftCode          string   `json:"swiftCode,omitempty"`
	PayeeCountryCode   string   `json:"payeeCountryCode,omitempty"`
	Address            string   `json:"address,omitempty"`
	PayeeAddressSecond string   `json:"payeeAddressSecond,omitempty"`
	MiddleSwiftCode    string   `json:"middleSwiftCode,omitempty"`
	BankName           string   `json:"bankName,omitempty"`
	BankCode           string   `json:"bankCode,omitempty"`
	BankAddress        string   `json:"bankAddress,omitempty"`
	AccountTypes       string   `json:"accountTypes,omitempty"`
}

type FiatAddressAddResp struct {
	AccountId string `json:"accountId"`
}

// FiatAddressDelReq is the body for POST /rest/v2.0/customer/fiat/withdrawal/del.
// See: https://opendocs.kun.global/docs/api/unbind-fiat-withdrawal-account
type FiatAddressDelReq struct {
	RequestNo string `json:"requestNo"`
	AccountId string `json:"accountId"`
	Currency  string `json:"currency"`
}

type RegionBalanceQueryReq struct {
	RequestNo    string `json:"requestNo"`
	Currency     string `json:"currency,omitempty"`
	CurrencyType string `json:"currencyType,omitempty"`
	RegionCode   string `json:"regionCode,omitempty"`
}

// AccountBalanceItem is a balance record returned by KUN account balance APIs.
// See: https://opendocs.kun.global/docs/api/get-account-balance
// See: https://opendocs.kun.global/docs/api/get-regional-account-balance
type AccountBalanceItem struct {
	Currency   string `json:"currency"`
	Balance    string `json:"balance"`
	RegionCode string `json:"regionCode,omitempty"`
}

// OutAccountBalanceQueryReq is the body for POST /rest/v2.0/trade/account/outAccount/query.
// See: https://opendocs.kun.global/docs/api/get-account-balance
type OutAccountBalanceQueryReq struct {
	RequestNo    string `json:"requestNo"`
	Currency     string `json:"currency"`
	CurrencyType string `json:"currencyType"`
}
