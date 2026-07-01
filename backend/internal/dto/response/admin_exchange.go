package response

import "time"

type AdminExchangeListItem struct {
	ID              uint64     `json:"id"`
	PlatformOrderID string     `json:"platform_order_id"`
	MerchantID      uint64     `json:"merchant_id"`
	MerchantEmail   string     `json:"merchant_email"`
	ExchangeType    string     `json:"exchange_type"`
	FromCurrency    string     `json:"from_currency"`
	ToCurrency      string     `json:"to_currency"`
	FromAmount      string     `json:"from_amount"`
	ToAmount        string     `json:"to_amount,omitempty"`
	ExchangeRate    string     `json:"exchange_rate,omitempty"`
	PlatformFee     string     `json:"platform_fee"`
	Status          string     `json:"status"`
	FailReason      string     `json:"fail_reason,omitempty"`
	KunOrderID      *string    `json:"kun_order_id,omitempty"`
	KunRequestNo    *string    `json:"kun_request_no,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

type AdminExchangeListResp struct {
	Exchanges []AdminExchangeListItem `json:"exchanges"`
	Total     int64                   `json:"total"`
}

type AdminExchangeSyncResp struct {
	OrderID       uint64 `json:"order_id"`
	Status        string `json:"status"`
	KunStatus     string `json:"kun_status"`
	Updated       bool   `json:"updated"`
	ToAmount      string `json:"to_amount,omitempty"`
	ExchangeRate  string `json:"exchange_rate,omitempty"`
	PlatformFee   string `json:"platform_fee"`
	FailReason    string `json:"fail_reason,omitempty"`
}
