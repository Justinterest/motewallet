package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet-withdrawal/backend/internal/model"
)

type FeeTemplateCryptoWithdrawalItemRepository interface {
	BatchReplace(ctx context.Context, tx *gorm.DB, templateID uint64, items []*model.FeeTemplateCryptoWithdrawalItem) error
	FindByTemplateID(ctx context.Context, templateID uint64) ([]*model.FeeTemplateCryptoWithdrawalItem, error)
}

type feeTemplateCryptoWithdrawalItemRepository struct {
	db *gorm.DB
}

func NewFeeTemplateCryptoWithdrawalItemRepository(db *gorm.DB) FeeTemplateCryptoWithdrawalItemRepository {
	return &feeTemplateCryptoWithdrawalItemRepository{db: db}
}

func (r *feeTemplateCryptoWithdrawalItemRepository) BatchReplace(ctx context.Context, tx *gorm.DB, templateID uint64, items []*model.FeeTemplateCryptoWithdrawalItem) error {
	if err := tx.WithContext(ctx).Where("fee_template_id = ?", templateID).Delete(&model.FeeTemplateCryptoWithdrawalItem{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&items).Error
}

func (r *feeTemplateCryptoWithdrawalItemRepository) FindByTemplateID(ctx context.Context, templateID uint64) ([]*model.FeeTemplateCryptoWithdrawalItem, error) {
	var items []*model.FeeTemplateCryptoWithdrawalItem
	err := r.db.WithContext(ctx).Where("fee_template_id = ?", templateID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
