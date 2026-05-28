package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type MerchantRepository interface {
	Create(ctx context.Context, merchant *model.Merchant) error
	FindByEmail(ctx context.Context, email string) (*model.Merchant, error)
	FindByID(ctx context.Context, id uint64) (*model.Merchant, error)
	Update(ctx context.Context, merchant *model.Merchant) error
	UpdateStatus(ctx context.Context, id uint64, status string) error
	ListWithPagination(ctx context.Context, page, pageSize int, status, kycStatus, search string) ([]*model.Merchant, int64, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
	FindByKunSubCustomerNo(ctx context.Context, subCustomerNo string) (*model.Merchant, error)
}

type merchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

func (r *merchantRepository) Create(ctx context.Context, merchant *model.Merchant) error {
	return r.db.WithContext(ctx).Create(merchant).Error
}

func (r *merchantRepository) FindByEmail(ctx context.Context, email string) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *merchantRepository) FindByID(ctx context.Context, id uint64) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *merchantRepository) Update(ctx context.Context, merchant *model.Merchant) error {
	return r.db.WithContext(ctx).Save(merchant).Error
}

func (r *merchantRepository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	return r.db.WithContext(ctx).Model(&model.Merchant{}).Where("id = ?", id).Update("status", status).Error
}

func (r *merchantRepository) ListWithPagination(ctx context.Context, page, pageSize int, status, kycStatus, search string) ([]*model.Merchant, int64, error) {
	var merchants []*model.Merchant
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Merchant{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if kycStatus != "" {
		query = query.Where("kyc_status = ?", kycStatus)
	}
	if search != "" {
		query = query.Where("email LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&merchants).Error; err != nil {
		return nil, 0, err
	}

	return merchants, total, nil
}

func (r *merchantRepository) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Merchant{}).Where("id = ?", id).Updates(fields).Error
}

func (r *merchantRepository) FindByKunSubCustomerNo(ctx context.Context, subCustomerNo string) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.WithContext(ctx).Where("kun_sub_customer_no = ?", subCustomerNo).First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}
