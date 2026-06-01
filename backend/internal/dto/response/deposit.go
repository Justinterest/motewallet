package response

import "time"

type DepositAddressResp struct {
	Address  string `json:"address"`
	Currency string `json:"currency"`
	Network  string `json:"network"`
}

type DepositOrderResp struct {
	ID        string    `json:"id"`
	Currency  string    `json:"currency"`
	Network   string    `json:"network"`
	Amount    string    `json:"amount"`
	TxHash    *string   `json:"tx_hash,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type DepositOrderListResp struct {
	Orders []DepositOrderResp `json:"orders"`
	Total  int64              `json:"total"`
}
