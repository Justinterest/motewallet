package dto

type WebhookEvent struct {
	EventID    string                 `json:"eventId"`
	EventType  string                 `json:"eventType"`
	EventTopic string                 `json:"eventTopic"`
	Timestamp  string                 `json:"timestamp"`
	Sign       string                 `json:"sign"`
	Data       map[string]interface{} `json:"data"`
}

type CustomerAuthData struct {
	SubCustomerNo string `json:"subCustomerNo"`
	AuthStatus    string `json:"authStatus"`
	FailReason    string `json:"failReason,omitempty"`
}

type CryptoDepositData struct {
	OrderId       string `json:"orderId"`
	SubCustomerNo string `json:"subCustomerNo"`
	Currency      string `json:"currency"`
	Chain         string `json:"chain"`
	OrderAmount   string `json:"orderAmount"`
	OrderCurrency string `json:"orderCurrency"`
	OrderStatus   string `json:"orderStatus"`
	TxId          string `json:"txId"`
	ToAddress     string `json:"toAddress"`
	FromAddress   string `json:"fromAddress"`
	FeeAmount     string `json:"feeAmount"`
	FeeCurrency   string `json:"feeCurrency"`
}

type CryptoWithdrawalData struct {
	OrderId       string `json:"orderId"`
	SubCustomerNo string `json:"subCustomerNo"`
	RequestNo     string `json:"requestNo"`
	Currency      string `json:"currency"`
	Chain         string `json:"chain"`
	Amount        string `json:"amount"`
	TxId          string `json:"txId"`
	OrderStatus   string `json:"orderStatus"`
	FeeAmount     string `json:"feeAmount"`
	FeeCurrency   string `json:"feeCurrency"`
}

type FiatWithdrawalData struct {
	OrderId       string `json:"orderId"`
	SubCustomerNo string `json:"subCustomerNo"`
	RequestNo     string `json:"requestNo"`
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	OrderStatus   string `json:"orderStatus"`
	FeeAmount     string `json:"feeAmount"`
	FeeCurrency   string `json:"feeCurrency"`
}

type ExchangeData struct {
	OrderId       string `json:"orderId"`
	SubCustomerNo string `json:"subCustomerNo"`
	RequestNo     string `json:"requestNo"`
	FromCurrency  string `json:"fromCurrency"`
	ToCurrency    string `json:"toCurrency"`
	FromAmount    string `json:"fromAmount"`
	ToAmount      string `json:"toAmount"`
	ExchangeRate  string `json:"exchangeRate"`
	OrderStatus   string `json:"orderStatus"`
	TradeFee      string `json:"tradeFee"`
	FeeCurrency   string `json:"feeCurrency"`
}
