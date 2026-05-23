package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet-withdrawal/backend/internal/model"
)

type FeeTemplateExchangeItemRepository interface {
	BatchReplace(ctx context.Context, tx *gorm.DB, templateID uint64, items []*model.FeeTemplateExchangeItem) error
	FindByTemplateID(ctx context.Context, templateID uint64) ([]*model.FeeTemplateExchangeItem, error)
}

type feeTemplateExchangeItemRepository struct {
	db *gorm.DB
}

func NewFeeTemplateExchangeItemRepository(db *gorm.DB) FeeTemplateExchangeItemRepository {
	return &feeTemplateExchangeItemRepository{db: db}
}

func (r *feeTemplateExchangeItemRepository) BatchReplace(ctx context.Context, tx *gorm.DB, templateID uint64, items []*model.FeeTemplateExchangeItem) error {
	if err := tx.WithContext(ctx).Where("fee_template_id = ?", templateID).Delete(&model.FeeTemplateExchangeItem{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&items).Error
}

func (r *feeTemplateExchangeItemRepository) FindByTemplateID(ctx context.Context, templateID uint64) ([]*model.FeeTemplateExchangeItem, error) {
	var items []*model.FeeTemplateExchangeItem
	err := r.db.WithContext(ctx).Where("fee_template_id = ?", templateID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
