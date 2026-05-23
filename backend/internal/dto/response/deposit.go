package response

import "time"

type DepositAddressResp struct {
	Address  string `json:"address"`
	Currency string `json:"currency"`
	Chain    string `json:"chain"`
}

type DepositOrderResp struct {
	ID        uint64    `json:"id"`
	Currency  string    `json:"currency"`
	Chain     string    `json:"chain"`
	Amount    string    `json:"amount"`
	TxID      *string   `json:"tx_id,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type DepositOrderListResp struct {
	Orders []DepositOrderResp `json:"orders"`
	Total  int64              `json:"total"`
}
