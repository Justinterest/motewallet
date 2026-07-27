package request

type ExchangeItemReq struct {
	FromCurrency   string `json:"from_currency" binding:"required"`
	ToCurrency     string `json:"to_currency" binding:"required"`
	FeeRate        string `json:"fee_rate" binding:"required"`
	MinFee         string `json:"min_fee" binding:"required"`
	MinFeeCurrency string `json:"min_fee_currency" binding:"required"`
}

type CryptoWithdrawalItemReq struct {
	Currency string `json:"currency" binding:"required"`
	Chain    string `json:"chain" binding:"required"`
	FeeRate  string `json:"fee_rate" binding:"required"`
	FixedFee string `json:"fixed_fee" binding:"required"`
}

type FiatWithdrawalItemReq struct {
	Currency     string `json:"currency" binding:"required"`
	TransferType string `json:"transfer_type" binding:"required"`
	FeeRate      string `json:"fee_rate" binding:"required"`
	FixedFee     string `json:"fixed_fee" binding:"required"`
}

type CreateFeeTemplateReq struct {
	Name                               string                    `json:"name" binding:"required"`
	Description                        *string                   `json:"description"`
	IsDefault                          bool                      `json:"is_default"`
	ExchangeFeeDeductionMethod         string                    `json:"exchange_fee_deduction_method"`
	CryptoWithdrawalFeeDeductionMethod string                    `json:"crypto_withdrawal_fee_deduction_method"`
	FiatWithdrawalFeeDeductionMethod   string                    `json:"fiat_withdrawal_fee_deduction_method"`
	ExchangeItems                      []ExchangeItemReq         `json:"exchange_items"`
	CryptoWithdrawalItems              []CryptoWithdrawalItemReq `json:"crypto_withdrawal_items"`
	FiatWithdrawalItems                []FiatWithdrawalItemReq   `json:"fiat_withdrawal_items"`
}

type UpdateFeeTemplateReq struct {
	Name                               string                    `json:"name" binding:"required"`
	Description                        *string                   `json:"description"`
	IsDefault                          bool                      `json:"is_default"`
	ExchangeFeeDeductionMethod         string                    `json:"exchange_fee_deduction_method"`
	CryptoWithdrawalFeeDeductionMethod string                    `json:"crypto_withdrawal_fee_deduction_method"`
	FiatWithdrawalFeeDeductionMethod   string                    `json:"fiat_withdrawal_fee_deduction_method"`
	ExchangeItems                      []ExchangeItemReq         `json:"exchange_items"`
	CryptoWithdrawalItems              []CryptoWithdrawalItemReq `json:"crypto_withdrawal_items"`
	FiatWithdrawalItems                []FiatWithdrawalItemReq   `json:"fiat_withdrawal_items"`
}
