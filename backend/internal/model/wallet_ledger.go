package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Wallet ledger entry types — one row per wallet mutation.
const (
	WalletLedgerCredit       = "CREDIT"
	WalletLedgerFreeze       = "FREEZE"
	WalletLedgerUnfreeze     = "UNFREEZE"
	WalletLedgerDeductFrozen = "DEDUCT_FROZEN"
)

// WalletLedger is an append-only record of a single wallet balance/frozen change.
type WalletLedger struct {
	ID                  uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID          uint64          `gorm:"column:merchant_id;not null;index:idx_wallet_ledger_merchant_currency_time" json:"merchant_id"`
	WalletID            uint64          `gorm:"column:wallet_id;not null;index:idx_wallet_ledger_wallet_time" json:"wallet_id"`
	AccountType         string          `gorm:"column:account_type;type:varchar(10);not null" json:"account_type"`
	Currency            string          `gorm:"column:currency;type:varchar(10);not null;index:idx_wallet_ledger_merchant_currency_time" json:"currency"`
	EntryType           string          `gorm:"column:entry_type;type:varchar(20);not null" json:"entry_type"`
	Amount              decimal.Decimal `gorm:"column:amount;type:decimal(28,8);not null" json:"amount"`
	BalanceBefore       decimal.Decimal `gorm:"column:balance_before;type:decimal(28,8);not null" json:"balance_before"`
	BalanceAfter        decimal.Decimal `gorm:"column:balance_after;type:decimal(28,8);not null" json:"balance_after"`
	FrozenBefore        decimal.Decimal `gorm:"column:frozen_before;type:decimal(28,8);not null" json:"frozen_before"`
	FrozenAfter         decimal.Decimal `gorm:"column:frozen_after;type:decimal(28,8);not null" json:"frozen_after"`
	TransactionRecordID *uint64         `gorm:"column:transaction_record_id;index:idx_wallet_ledger_txn" json:"transaction_record_id"`
	BizType             *string         `gorm:"column:biz_type;type:varchar(20);index:idx_wallet_ledger_biz_type" json:"biz_type"`
	Remark              *string         `gorm:"column:remark;type:varchar(500)" json:"remark"`
	CreatedAt           time.Time       `gorm:"column:created_at;type:datetime(3);not null;autoCreateTime;index:idx_wallet_ledger_merchant_currency_time;index:idx_wallet_ledger_wallet_time" json:"created_at"`
}

func (WalletLedger) TableName() string {
	return "wallet_ledger"
}
