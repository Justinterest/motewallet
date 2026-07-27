package response

import "time"

type WithdrawalOrderResp struct {
	ID                 uint64    `json:"id"`
	Type               string    `json:"type"`
	Currency           string    `json:"currency"`
	Chain              *string   `json:"chain,omitempty"`
	Amount             string    `json:"amount"`
	PlatformFee        string    `json:"platform_fee"`
	FeeDeductionMethod string    `json:"fee_deduction_method"`
	NetAmount          string    `json:"net_amount"`
	Status             string    `json:"status"`
	ReviewStatus       string    `json:"review_status"`
	ToAddress          *string   `json:"to_address,omitempty"`
	TxID               *string   `json:"tx_id,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

type WithdrawalOrderListResp struct {
	Orders []WithdrawalOrderResp `json:"orders"`
	Total  int64                 `json:"total"`
}

type WithdrawalFeePreviewResp struct {
	Currency           string `json:"currency"`
	Amount             string `json:"amount"`
	PlatformFee        string `json:"platform_fee"`
	FeeDeductionMethod string `json:"fee_deduction_method"`
	TotalDeduction     string `json:"total_deduction"`
	NetAmount          string `json:"net_amount"`
}
