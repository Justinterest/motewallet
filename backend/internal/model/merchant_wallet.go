package model

import "github.com/shopspring/decimal"

type MerchantWallet struct {
	BaseModel
	MerchantID    uint64          `gorm:"column:merchant_id;not null;uniqueIndex:uk_merchant_wallets_account;index:idx_merchant_wallets_merchant_id" json:"merchant_id"`
	AccountType   string          `gorm:"column:account_type;type:varchar(10);not null;uniqueIndex:uk_merchant_wallets_account" json:"account_type"`
	Currency      string          `gorm:"column:currency;type:varchar(10);not null;uniqueIndex:uk_merchant_wallets_account" json:"currency"`
	Balance       decimal.Decimal `gorm:"column:balance;type:decimal(28,8);not null;default:0" json:"balance"`
	FrozenBalance decimal.Decimal `gorm:"column:frozen_balance;type:decimal(28,8);not null;default:0" json:"frozen_balance"`
	Version       uint64          `gorm:"column:version;not null;default:0" json:"version"`
}

func (MerchantWallet) TableName() string {
	return "merchant_wallets"
}
