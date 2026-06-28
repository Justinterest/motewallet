package dto

type CryptoWithdrawalReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	RequestNo     string `json:"requestNo"`
	Currency      string `json:"currency"`
	Chain         string `json:"chain"`
	Amount        string `json:"amount"`
	ToAddress     string `json:"toAddress"`
	RegionCode    string `json:"regionCode"`
}

type CryptoWithdrawalResp struct {
	OrderId string `json:"orderId"`
}

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
	SubCustomerNo string `json:"subCustomerNo"`
	RequestNo     string `json:"requestNo"`
	FromCurrency  string `json:"fromCurrency"`
	ToCurrency    string `json:"toCurrency"`
	FromAmount    string `json:"fromAmount"`
}

type InnerMatchCreateResp struct {
	OrderId string `json:"orderId"`
}

type InnerMatchQueryReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	OrderId       string `json:"orderId"`
}

type InnerMatchQueryResp struct {
	OrderId      string `json:"orderId"`
	OrderStatus  string `json:"orderStatus"`
	FromCurrency string `json:"fromCurrency"`
	ToCurrency   string `json:"toCurrency"`
	FromAmount   string `json:"fromAmount"`
	ToAmount     string `json:"toAmount"`
	ExchangeRate string `json:"exchangeRate"`
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
}
