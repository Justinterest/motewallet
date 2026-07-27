package model

type FeeTemplate struct {
	SoftDeleteModel
	Name                               string  `gorm:"column:name;type:varchar(64);not null" json:"name"`
	Description                        *string `gorm:"column:description;type:varchar(255)" json:"description,omitempty"`
	IsDefault                          bool    `gorm:"column:is_default;type:tinyint(1);not null;default:0;index:idx_fee_templates_is_default" json:"is_default"`
	ExchangeFeeDeductionMethod         string  `gorm:"column:exchange_fee_deduction_method;type:varchar(20);not null;default:WALLET" json:"exchange_fee_deduction_method"`
	CryptoWithdrawalFeeDeductionMethod string  `gorm:"column:crypto_withdrawal_fee_deduction_method;type:varchar(20);not null;default:WALLET" json:"crypto_withdrawal_fee_deduction_method"`
	FiatWithdrawalFeeDeductionMethod   string  `gorm:"column:fiat_withdrawal_fee_deduction_method;type:varchar(20);not null;default:WALLET" json:"fiat_withdrawal_fee_deduction_method"`
}

func (FeeTemplate) TableName() string {
	return "fee_templates"
}
