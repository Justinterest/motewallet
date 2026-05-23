package model

import "github.com/shopspring/decimal"

type FeeTemplateFiatWithdrawalItem struct {
	BaseModel
	FeeTemplateID uint64          `gorm:"column:fee_template_id;not null;uniqueIndex:uk_fiat_withdrawal_items_template_type" json:"fee_template_id"`
	Currency      string          `gorm:"column:currency;type:varchar(10);not null;uniqueIndex:uk_fiat_withdrawal_items_template_type" json:"currency"`
	TransferType  string          `gorm:"column:transfer_type;type:varchar(10);not null;uniqueIndex:uk_fiat_withdrawal_items_template_type" json:"transfer_type"`
	FeeRate       decimal.Decimal `gorm:"column:fee_rate;type:decimal(10,6);not null;default:0" json:"fee_rate"`
	FixedFee      decimal.Decimal `gorm:"column:fixed_fee;type:decimal(28,8);not null;default:0" json:"fixed_fee"`
}

func (FeeTemplateFiatWithdrawalItem) TableName() string {
	return "fee_template_fiat_withdrawal_items"
}
