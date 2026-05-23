package model

import "github.com/shopspring/decimal"

type FeeTemplateCryptoWithdrawalItem struct {
	BaseModel
	FeeTemplateID uint64          `gorm:"column:fee_template_id;not null;uniqueIndex:uk_crypto_withdrawal_items_template_chain" json:"fee_template_id"`
	Currency      string          `gorm:"column:currency;type:varchar(10);not null;uniqueIndex:uk_crypto_withdrawal_items_template_chain" json:"currency"`
	Chain         string          `gorm:"column:chain;type:varchar(20);not null;uniqueIndex:uk_crypto_withdrawal_items_template_chain" json:"chain"`
	FeeRate       decimal.Decimal `gorm:"column:fee_rate;type:decimal(10,6);not null;default:0" json:"fee_rate"`
	FixedFee      decimal.Decimal `gorm:"column:fixed_fee;type:decimal(28,8);not null;default:0" json:"fixed_fee"`
}

func (FeeTemplateCryptoWithdrawalItem) TableName() string {
	return "fee_template_crypto_withdrawal_items"
}
