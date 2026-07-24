package repository

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"motewallet/internal/model"
)

var ErrOptimisticLock = errors.New("optimistic lock: version mismatch")

type MerchantWalletRepository interface {
	CreateBatch(ctx context.Context, wallets []*model.MerchantWallet) error
	FindByMerchantID(ctx context.Context, merchantID uint64) ([]*model.MerchantWallet, error)
	FindByMerchantAndAccount(ctx context.Context, merchantID uint64, accountType, currency string) (*model.MerchantWallet, error)
	FindByMerchantAccountCurrency(ctx context.Context, merchantID uint64, accountType, currency string) (*model.MerchantWallet, error)
	FindByMerchantAccountCurrencyWithDB(ctx context.Context, db *gorm.DB, merchantID uint64, accountType, currency string) (*model.MerchantWallet, error)
	// FindOrCreateWithDB returns the wallet row, creating a zero-balance wallet if missing.
	FindOrCreateWithDB(ctx context.Context, db *gorm.DB, merchantID uint64, accountType, currency string) (*model.MerchantWallet, error)
	UpdateBalanceWithVersion(ctx context.Context, db *gorm.DB, wallet *model.MerchantWallet) error
}

type merchantWalletRepository struct {
	db *gorm.DB
}

func NewMerchantWalletRepository(db *gorm.DB) MerchantWalletRepository {
	return &merchantWalletRepository{db: db}
}

func (r *merchantWalletRepository) CreateBatch(ctx context.Context, wallets []*model.MerchantWallet) error {
	return r.db.WithContext(ctx).Create(&wallets).Error
}

func (r *merchantWalletRepository) FindByMerchantID(ctx context.Context, merchantID uint64) ([]*model.MerchantWallet, error) {
	var wallets []*model.MerchantWallet
	err := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *merchantWalletRepository) FindByMerchantAndAccount(ctx context.Context, merchantID uint64, accountType, currency string) (*model.MerchantWallet, error) {
	var wallet model.MerchantWallet
	err := r.db.WithContext(ctx).Where("merchant_id = ? AND account_type = ? AND currency = ?", merchantID, accountType, currency).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *merchantWalletRepository) FindByMerchantAccountCurrency(ctx context.Context, merchantID uint64, accountType, currency string) (*model.MerchantWallet, error) {
	var wallet model.MerchantWallet
	err := r.db.WithContext(ctx).Where("merchant_id = ? AND account_type = ? AND currency = ?", merchantID, accountType, currency).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *merchantWalletRepository) FindByMerchantAccountCurrencyWithDB(ctx context.Context, db *gorm.DB, merchantID uint64, accountType, currency string) (*model.MerchantWallet, error) {
	var wallet model.MerchantWallet
	err := db.WithContext(ctx).Where("merchant_id = ? AND account_type = ? AND currency = ?", merchantID, accountType, currency).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *merchantWalletRepository) FindOrCreateWithDB(ctx context.Context, db *gorm.DB, merchantID uint64, accountType, currency string) (*model.MerchantWallet, error) {
	wallet, err := r.FindByMerchantAccountCurrencyWithDB(ctx, db, merchantID, accountType, currency)
	if err == nil {
		return wallet, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	created := &model.MerchantWallet{
		MerchantID:    merchantID,
		AccountType:   accountType,
		Currency:      currency,
		Balance:       decimal.Zero,
		FrozenBalance: decimal.Zero,
		Version:       0,
	}
	// Ignore duplicate-key races from concurrent ensure calls; re-read below.
	if err := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(created).Error; err != nil {
		return nil, err
	}

	return r.FindByMerchantAccountCurrencyWithDB(ctx, db, merchantID, accountType, currency)
}

func (r *merchantWalletRepository) UpdateBalanceWithVersion(ctx context.Context, db *gorm.DB, wallet *model.MerchantWallet) error {
	result := db.WithContext(ctx).Model(&model.MerchantWallet{}).
		Where("id = ? AND version = ?", wallet.ID, wallet.Version).
		Updates(map[string]interface{}{
			"balance":        wallet.Balance,
			"frozen_balance": wallet.FrozenBalance,
			"version":        gorm.Expr("version + 1"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLock
	}
	return nil
}
