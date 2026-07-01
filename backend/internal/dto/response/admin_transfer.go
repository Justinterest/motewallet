package response

import "time"

type AdminTransferListItem struct {
	ID              uint64    `json:"id"`
	PlatformOrderID string    `json:"platform_order_id"`
	MerchantID      uint64    `json:"merchant_id"`
	MerchantEmail   string    `json:"merchant_email"`
	FromAccountType string    `json:"from_account_type"`
	ToAccountType   string    `json:"to_account_type"`
	Currency        string    `json:"currency"`
	Amount          string    `json:"amount"`
	Status          string    `json:"status"`
	KunOrderID      *string   `json:"kun_order_id,omitempty"`
	KunRequestNo    *string   `json:"kun_request_no,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type AdminTransferListResp struct {
	Transfers []AdminTransferListItem `json:"transfers"`
	Total     int64                   `json:"total"`
}

type AdminTransferSyncResp struct {
	OrderID   uint64 `json:"order_id"`
	Status    string `json:"status"`
	KunStatus string `json:"kun_status"`
	Updated   bool   `json:"updated"`
}
