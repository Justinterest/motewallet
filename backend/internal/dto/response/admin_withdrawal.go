package response

import "time"

type AdminWithdrawalListItem struct {
	ID              uint64     `json:"id"`
	PlatformOrderID string     `json:"platform_order_id"`
	MerchantID      uint64     `json:"merchant_id"`
	MerchantEmail   string     `json:"merchant_email"`
	Type            string     `json:"type"`
	Currency        string     `json:"currency"`
	Network         string     `json:"network,omitempty"`
	Amount          string     `json:"amount"`
	PlatformFee     string     `json:"platform_fee"`
	Status          string     `json:"status"`
	ReviewStatus    string     `json:"review_status"`
	ToAddress       *string    `json:"to_address,omitempty"`
	TxID            *string    `json:"tx_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

type AdminWithdrawalListResp struct {
	Withdrawals []AdminWithdrawalListItem `json:"withdrawals"`
	Total       int64                     `json:"total"`
}
