package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type FeeTemplateFiatWithdrawalItemRepository interface {
	BatchReplace(ctx context.Context, tx *gorm.DB, templateID uint64, items []*model.FeeTemplateFiatWithdrawalItem) error
	FindByTemplateID(ctx context.Context, templateID uint64) ([]*model.FeeTemplateFiatWithdrawalItem, error)
}

type feeTemplateFiatWithdrawalItemRepository struct {
	db *gorm.DB
}

func NewFeeTemplateFiatWithdrawalItemRepository(db *gorm.DB) FeeTemplateFiatWithdrawalItemRepository {
	return &feeTemplateFiatWithdrawalItemRepository{db: db}
}

func (r *feeTemplateFiatWithdrawalItemRepository) BatchReplace(ctx context.Context, tx *gorm.DB, templateID uint64, items []*model.FeeTemplateFiatWithdrawalItem) error {
	if err := tx.WithContext(ctx).Where("fee_template_id = ?", templateID).Delete(&model.FeeTemplateFiatWithdrawalItem{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&items).Error
}

func (r *feeTemplateFiatWithdrawalItemRepository) FindByTemplateID(ctx context.Context, templateID uint64) ([]*model.FeeTemplateFiatWithdrawalItem, error) {
	var items []*model.FeeTemplateFiatWithdrawalItem
	err := r.db.WithContext(ctx).Where("fee_template_id = ?", templateID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
