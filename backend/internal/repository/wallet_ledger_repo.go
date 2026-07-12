package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type WalletLedgerListFilter struct {
	AccountType string
	Currency    string
	BizType     string
	EntryType   string
	Page        int
	PageSize    int
}

type WalletLedgerListItem struct {
	model.WalletLedger
	PlatformOrderID *string `gorm:"column:platform_order_id"`
}

type WalletLedgerRepository interface {
	Create(ctx context.Context, db *gorm.DB, entry *model.WalletLedger) error
	ListByMerchant(ctx context.Context, merchantID uint64, filter WalletLedgerListFilter) ([]WalletLedgerListItem, int64, error)
}

type walletLedgerRepository struct {
	db *gorm.DB
}

func NewWalletLedgerRepository(db *gorm.DB) WalletLedgerRepository {
	return &walletLedgerRepository{db: db}
}

func (r *walletLedgerRepository) Create(ctx context.Context, db *gorm.DB, entry *model.WalletLedger) error {
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Create(entry).Error
}

func (r *walletLedgerRepository) ListByMerchant(ctx context.Context, merchantID uint64, filter WalletLedgerListFilter) ([]WalletLedgerListItem, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := r.db.WithContext(ctx).Table("wallet_ledger").Where("wallet_ledger.merchant_id = ?", merchantID)
	if filter.AccountType != "" {
		query = query.Where("wallet_ledger.account_type = ?", filter.AccountType)
	}
	if filter.Currency != "" {
		query = query.Where("wallet_ledger.currency = ?", filter.Currency)
	}
	if filter.BizType != "" {
		query = query.Where("wallet_ledger.biz_type = ?", filter.BizType)
	}
	if filter.EntryType != "" {
		query = query.Where("wallet_ledger.entry_type = ?", filter.EntryType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []WalletLedgerListItem
	err := query.
		Select("wallet_ledger.*, transaction_records.platform_order_id AS platform_order_id").
		Joins("LEFT JOIN transaction_records ON transaction_records.id = wallet_ledger.transaction_record_id").
		Order("wallet_ledger.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
