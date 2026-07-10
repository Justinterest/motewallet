package dto

type CryptoWithdrawalReq struct {
	RequestNo  string `json:"requestNo"`
	Amount     string `json:"amount"`
	Chain      string `json:"chain"`
	Currency   string `json:"currency"`
	Address    string `json:"address"`
	RegionCode string `json:"regionCode"`
}

// CryptoWithdrawalResp order id is returned as a plain string in data.
type CryptoWithdrawalResp string

// FiatWithdrawalReq is the body for POST /rest/v2.0/trade/fiat/withdrawal.
// See: https://opendocs.kun.global/docs/api/fiat-withdrawal
type FiatWithdrawalReq struct {
	RequestNo  string `json:"requestNo"`
	AccountId  string `json:"accountId"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	FeeMethod  string `json:"feeMethod"`
	PoboType   string `json:"poboType"`
	Postscript string `json:"postscript"`
	Purpose    string `json:"purpose"`
}

type FiatWithdrawalResp struct {
	OrderId string `json:"orderId"`
}

// ExchangeQuoteReq is the body for POST /rest/v2.0/trade/exchange/quote/request.
// See: https://opendocs.kun.global/docs/api/quote-request
type ExchangeQuoteReq struct {
	RequestNo      string `json:"requestNo"`
	Amount         string `json:"amount"`
	QuoteCurrency  string `json:"quoteCurrency"`
	QuotedCurrency string `json:"quotedCurrency"`
}

type ExchangeQuoteResp struct {
	QuoteId      string `json:"quoteId"`
	FromCurrency string `json:"fromCurrency"`
	ToCurrency   string `json:"toCurrency"`
	FromAmount   string `json:"fromAmount"`
	ToAmount     string `json:"toAmount"`
	ExchangeRate string `json:"exchangeRate"`
	TradeFee     string `json:"tradeFee"`
	FeeCurrency  string `json:"feeCurrency"`
	ExpireTime   int64  `json:"expireTime"`
}

// ExchangeOrderReq is the body for POST /rest/v2.0/trade/exchange/order.
type ExchangeOrderReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	RequestNo     string `json:"requestNo"`
	QuoteId       string `json:"quoteId"`
	FromCurrency  string `json:"fromCurrency"`
	ToCurrency    string `json:"toCurrency"`
	FromAmount    string `json:"fromAmount"`
}

type ExchangeOrderResp struct {
	OrderId string `json:"orderId"`
}

type ExchangeOrderQueryReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	OrderId       string `json:"orderId"`
}

// ExchangeOrderQueryResp is the response for POST /rest/v2.0/trade/exchange/order/query.
type ExchangeOrderQueryResp struct {
	OrderId      string `json:"orderId"`
	OrderStatus  string `json:"orderStatus"`
	FromCurrency string `json:"fromCurrency"`
	ToCurrency   string `json:"toCurrency"`
	FromAmount   string `json:"fromAmount"`
	ToAmount     string `json:"toAmount"`
	ExchangeRate string `json:"exchangeRate"`
	TradeFee     string `json:"tradeFee"`
	FeeCurrency  string `json:"feeCurrency"`
}

type InnerMatchCreateReq struct {
	RequestNo    string `json:"requestNo"`
	FromCurrency string `json:"fromCurrency"`
	OrderAmount  string `json:"orderAmount"`
	ToCurrency   string `json:"toCurrency"`
	AutoTransfer string `json:"autoTransfer,omitempty"`
}

type InnerMatchCreateResp struct {
	OrderId string `json:"orderId"`
}

type InnerMatchQueryReq struct {
	RequestNo string `json:"requestNo"`
	OrderId   string `json:"orderId"`
}

type InnerMatchQueryResp struct {
	OrderId           string `json:"orderId"`
	OrderStatus       string `json:"orderStatus"`
	FromCurrency      string `json:"fromCurrency"`
	ToCurrency        string `json:"toCurrency"`
	OrderAmount       string `json:"orderAmount"`
	OrderCurrency     string `json:"orderCurrency"`
	ReceiveAmount     string `json:"receiveAmount"`
	ReceiveCurrency   string `json:"receiveCurrency"`
	TradeFee          string `json:"tradeFee"`
	TradeFeeCurrency  string `json:"tradeFeeCurrency"`
	ExchangeRate      string `json:"exchangeRate"`
	RejectReason      string `json:"rejectReason"`
	CompleteTime        string `json:"completeTime"`
}

// FundTransferReq is the body for POST /rest/v2.0/user/fund/transfer.
// See: https://opendocs.kun.global/docs/api/create-transfer
type FundTransferReq struct {
	RequestNo string `json:"requestNo"`
	FromAcc   string `json:"fromAcc"`
	ToAcc     string `json:"toAcc"`
	Currency  string `json:"currency"`
	Amount    string `json:"amount"`
}

type FundTransferResp struct {
	OrderId     string `json:"orderId"`
	OrderStatus string `json:"orderStatus"`
	Status      string `json:"status"`
}

func (r FundTransferResp) ResolvedStatus() string {
	if r.OrderStatus != "" {
		return r.OrderStatus
	}
	return r.Status
}

// FundTransferListReq is the body for POST /rest/v2.0/query/fund/transfer/list.
// See: https://opendocs.kun.global/docs/api/list-transfer-records
type FundTransferListReq struct {
	RequestNo         string `json:"requestNo"`
	OriginalRequestNo string `json:"originalRequestNo,omitempty"`
	OrderId           string `json:"orderId,omitempty"`
	StartTime         string `json:"startTime,omitempty"`
	EndTime           string `json:"endTime,omitempty"`
	RegionCode        string `json:"regionCode,omitempty"`
	PageNo            int    `json:"pageNo,omitempty"`
	PageSize          int    `json:"pageSize,omitempty"`
}

type FundTransferListItem struct {
	OrderId      string `json:"orderId"`
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
	FromAcc      string `json:"fromAcc"`
	ToAcc        string `json:"toAcc"`
	Currency     string `json:"currency"`
	Amount       string `json:"amount"`
	Status       string `json:"status"`
	TransferTime string `json:"transferTime"`
	RegionCode   string `json:"regionCode"`
}

type FundTransferListResp struct {
	PageCount   string                 `json:"pageCount"`
	PageNo      string                 `json:"pageNo"`
	PageSize    string                 `json:"pageSize"`
	RecordCount string                 `json:"recordCount"`
	Records     []FundTransferListItem `json:"records"`
}
