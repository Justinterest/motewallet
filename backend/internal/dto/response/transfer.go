package response

import "time"

type TransferOrderResp struct {
	ID              uint64    `json:"id"`
	FromAccountType string    `json:"from_account_type"`
	ToAccountType   string    `json:"to_account_type"`
	Currency        string    `json:"currency"`
	Amount          string    `json:"amount"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type TransferOrderListResp struct {
	Orders []TransferOrderResp `json:"orders"`
	Total  int64               `json:"total"`
}
