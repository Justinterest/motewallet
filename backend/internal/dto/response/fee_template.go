package response

import "time"

type ExchangeItemResp struct {
	ID             uint64 `json:"id"`
	FromCurrency   string `json:"from_currency"`
	ToCurrency     string `json:"to_currency"`
	FeeRate        string `json:"fee_rate"`
	MinFee         string `json:"min_fee"`
	MinFeeCurrency string `json:"min_fee_currency"`
}

type CryptoWithdrawalItemResp struct {
	ID       uint64 `json:"id"`
	Currency string `json:"currency"`
	Chain    string `json:"chain"`
	FeeRate  string `json:"fee_rate"`
	FixedFee string `json:"fixed_fee"`
}

type FiatWithdrawalItemResp struct {
	ID           uint64 `json:"id"`
	Currency     string `json:"currency"`
	TransferType string `json:"transfer_type"`
	FeeRate      string `json:"fee_rate"`
	FixedFee     string `json:"fixed_fee"`
}

type FeeTemplateListItem struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FeeTemplateListResp struct {
	Templates []FeeTemplateListItem `json:"templates"`
}

type FeeTemplateDetailResp struct {
	ID                                 uint64                     `json:"id"`
	Name                               string                     `json:"name"`
	Description                        *string                    `json:"description,omitempty"`
	IsDefault                          bool                       `json:"is_default"`
	ExchangeFeeDeductionMethod         string                     `json:"exchange_fee_deduction_method"`
	CryptoWithdrawalFeeDeductionMethod string                     `json:"crypto_withdrawal_fee_deduction_method"`
	FiatWithdrawalFeeDeductionMethod   string                     `json:"fiat_withdrawal_fee_deduction_method"`
	ExchangeItems                      []ExchangeItemResp         `json:"exchange_items"`
	CryptoWithdrawalItems              []CryptoWithdrawalItemResp `json:"crypto_withdrawal_items"`
	FiatWithdrawalItems                []FiatWithdrawalItemResp   `json:"fiat_withdrawal_items"`
	CreatedAt                          time.Time                  `json:"created_at"`
	UpdatedAt                          time.Time                  `json:"updated_at"`
}
