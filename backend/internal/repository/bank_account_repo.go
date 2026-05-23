package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet-withdrawal/backend/internal/model"
)

type BankAccountRepository interface {
	Create(ctx context.Context, account *model.BankAccount) error
	FindByID(ctx context.Context, id uint64) (*model.BankAccount, error)
	ListByMerchant(ctx context.Context, merchantID uint64) ([]*model.BankAccount, error)
	Delete(ctx context.Context, id uint64) error
}

type bankAccountRepository struct {
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) BankAccountRepository {
	return &bankAccountRepository{db: db}
}

func (r *bankAccountRepository) Create(ctx context.Context, account *model.BankAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *bankAccountRepository) FindByID(ctx context.Context, id uint64) (*model.BankAccount, error) {
	var account model.BankAccount
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *bankAccountRepository) ListByMerchant(ctx context.Context, merchantID uint64) ([]*model.BankAccount, error) {
	var accounts []*model.BankAccount
	err := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).Order("id DESC").Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *bankAccountRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BankAccount{}, id).Error
}
