package response

import "time"

type ExchangePreviewResp struct {
	FromCurrency   string `json:"from_currency"`
	ToCurrency     string `json:"to_currency"`
	FromAmount     string `json:"from_amount"`
	ToAmount       string `json:"to_amount"`
	ExchangeRate   string `json:"exchange_rate"`
	PlatformFee    string `json:"platform_fee"`
	FeeCurrency    string `json:"fee_currency"`
	NetToAmount    string `json:"net_to_amount"`
	TotalDeduction string `json:"total_deduction"`
}

type ExchangeOrderResp struct {
	ID           uint64    `json:"id"`
	ExchangeType string    `json:"exchange_type"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	FromAmount   string    `json:"from_amount"`
	ToAmount     string    `json:"to_amount"`
	ExchangeRate string    `json:"exchange_rate"`
	PlatformFee  string    `json:"platform_fee"`
	Status       string    `json:"status"`
	FailReason   string    `json:"fail_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type ExchangeOrderListResp struct {
	Orders []ExchangeOrderResp `json:"orders"`
	Total  int64               `json:"total"`
}
