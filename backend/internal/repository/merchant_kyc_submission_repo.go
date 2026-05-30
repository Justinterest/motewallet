package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type MerchantKycSubmissionRepository interface {
	Create(ctx context.Context, submission *model.MerchantKycSubmission) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
	FindLatestByMerchantID(ctx context.Context, merchantID uint64) (*model.MerchantKycSubmission, error)
}

type merchantKycSubmissionRepository struct {
	db *gorm.DB
}

func NewMerchantKycSubmissionRepository(db *gorm.DB) MerchantKycSubmissionRepository {
	return &merchantKycSubmissionRepository{db: db}
}

func (r *merchantKycSubmissionRepository) Create(ctx context.Context, submission *model.MerchantKycSubmission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r *merchantKycSubmissionRepository) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.MerchantKycSubmission{}).Where("id = ?", id).Updates(fields).Error
}

func (r *merchantKycSubmissionRepository) FindLatestByMerchantID(ctx context.Context, merchantID uint64) (*model.MerchantKycSubmission, error) {
	var submission model.MerchantKycSubmission
	err := r.db.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		Order("id DESC").
		First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}
