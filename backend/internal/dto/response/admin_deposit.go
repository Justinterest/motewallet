package response

import "time"

type AdminDepositListItem struct {
	ID              uint64     `json:"id"`
	PlatformOrderID string     `json:"platform_order_id"`
	MerchantID      uint64     `json:"merchant_id"`
	MerchantEmail   string     `json:"merchant_email"`
	Currency        string     `json:"currency"`
	Network         string     `json:"network"`
	Amount          string     `json:"amount"`
	TxHash          *string    `json:"tx_hash,omitempty"`
	ToAddress       string     `json:"to_address"`
	FromAddress     *string    `json:"from_address,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

type AdminDepositListResp struct {
	Deposits []AdminDepositListItem `json:"deposits"`
	Total    int64                  `json:"total"`
}
