package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type FeeTemplateRepository interface {
	Create(ctx context.Context, template *model.FeeTemplate) error
	FindByID(ctx context.Context, id uint64) (*model.FeeTemplate, error)
	FindDefault(ctx context.Context) (*model.FeeTemplate, error)
	FindAll(ctx context.Context) ([]*model.FeeTemplate, error)
	Update(ctx context.Context, template *model.FeeTemplate) error
	Delete(ctx context.Context, id uint64) error
	SetDefault(ctx context.Context, id uint64) error
	IsReferencedByMerchant(ctx context.Context, id uint64) (bool, error)
}

type feeTemplateRepository struct {
	db *gorm.DB
}

func NewFeeTemplateRepository(db *gorm.DB) FeeTemplateRepository {
	return &feeTemplateRepository{db: db}
}

func (r *feeTemplateRepository) Create(ctx context.Context, template *model.FeeTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *feeTemplateRepository) FindByID(ctx context.Context, id uint64) (*model.FeeTemplate, error) {
	var template model.FeeTemplate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *feeTemplateRepository) FindDefault(ctx context.Context) (*model.FeeTemplate, error) {
	var template model.FeeTemplate
	err := r.db.WithContext(ctx).Where("is_default = ?", true).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *feeTemplateRepository) FindAll(ctx context.Context) ([]*model.FeeTemplate, error) {
	var templates []*model.FeeTemplate
	err := r.db.WithContext(ctx).Order("id DESC").Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *feeTemplateRepository) Update(ctx context.Context, template *model.FeeTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *feeTemplateRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.FeeTemplate{}, id).Error
}

func (r *feeTemplateRepository) SetDefault(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.FeeTemplate{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&model.FeeTemplate{}).Where("id = ?", id).Update("is_default", true).Error
	})
}

func (r *feeTemplateRepository) IsReferencedByMerchant(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Merchant{}).Where("fee_template_id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
