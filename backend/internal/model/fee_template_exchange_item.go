package model

import "github.com/shopspring/decimal"

type FeeTemplateExchangeItem struct {
	BaseModel
	FeeTemplateID  uint64          `gorm:"column:fee_template_id;not null;uniqueIndex:uk_exchange_items_template_pair" json:"fee_template_id"`
	FromCurrency   string          `gorm:"column:from_currency;type:varchar(10);not null;uniqueIndex:uk_exchange_items_template_pair" json:"from_currency"`
	ToCurrency     string          `gorm:"column:to_currency;type:varchar(10);not null;uniqueIndex:uk_exchange_items_template_pair" json:"to_currency"`
	FeeRate        decimal.Decimal `gorm:"column:fee_rate;type:decimal(10,6);not null;default:0" json:"fee_rate"`
	MinFee         decimal.Decimal `gorm:"column:min_fee;type:decimal(28,8);not null;default:0" json:"min_fee"`
	MinFeeCurrency string          `gorm:"column:min_fee_currency;type:varchar(10);not null" json:"min_fee_currency"`
}

func (FeeTemplateExchangeItem) TableName() string {
	return "fee_template_exchange_items"
}
