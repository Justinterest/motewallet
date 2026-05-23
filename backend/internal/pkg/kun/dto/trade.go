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

type FiatWithdrawalReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	RequestNo     string `json:"requestNo"`
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	BankAccountId string `json:"bankAccountId"`
	RegionCode    string `json:"regionCode"`
}

type FiatWithdrawalResp struct {
	OrderId string `json:"orderId"`
}

type ExchangeQuoteReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	FromCurrency  string `json:"fromCurrency"`
	ToCurrency    string `json:"toCurrency"`
	FromAmount    string `json:"fromAmount,omitempty"`
	ToAmount      string `json:"toAmount,omitempty"`
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

type FundTransferReq struct {
	SubCustomerNo   string `json:"subCustomerNo"`
	RequestNo       string `json:"requestNo"`
	Currency        string `json:"currency"`
	Amount          string `json:"amount"`
	FromAccountType string `json:"fromAccountType"`
	ToAccountType   string `json:"toAccountType"`
	RegionCode      string `json:"regionCode"`
}

type FundTransferResp struct {
	OrderId string `json:"orderId"`
}
