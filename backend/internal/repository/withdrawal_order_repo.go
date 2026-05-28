package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type WithdrawalOrderRepository interface {
	Create(ctx context.Context, order *model.WithdrawalOrder) error
	CreateWithDB(ctx context.Context, db *gorm.DB, order *model.WithdrawalOrder) error
	FindByID(ctx context.Context, id uint64) (*model.WithdrawalOrder, error)
	ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.WithdrawalOrder, int64, error)
	ListPendingReview(ctx context.Context, page, pageSize int) ([]*model.WithdrawalOrder, int64, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
	FindByKunRequestNo(ctx context.Context, requestNo string) (*model.WithdrawalOrder, error)
}

type withdrawalOrderRepository struct {
	db *gorm.DB
}

func NewWithdrawalOrderRepository(db *gorm.DB) WithdrawalOrderRepository {
	return &withdrawalOrderRepository{db: db}
}

func (r *withdrawalOrderRepository) Create(ctx context.Context, order *model.WithdrawalOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *withdrawalOrderRepository) CreateWithDB(ctx context.Context, db *gorm.DB, order *model.WithdrawalOrder) error {
	return db.WithContext(ctx).Create(order).Error
}

func (r *withdrawalOrderRepository) FindByID(ctx context.Context, id uint64) (*model.WithdrawalOrder, error) {
	var order model.WithdrawalOrder
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *withdrawalOrderRepository) ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.WithdrawalOrder, int64, error) {
	var orders []*model.WithdrawalOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&model.WithdrawalOrder{}).Where("merchant_id = ?", merchantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *withdrawalOrderRepository) ListPendingReview(ctx context.Context, page, pageSize int) ([]*model.WithdrawalOrder, int64, error) {
	var orders []*model.WithdrawalOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&model.WithdrawalOrder{}).Where("review_status = ?", "PENDING_REVIEW")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *withdrawalOrderRepository) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.WithdrawalOrder{}).Where("id = ?", id).Updates(fields).Error
}

func (r *withdrawalOrderRepository) FindByKunRequestNo(ctx context.Context, requestNo string) (*model.WithdrawalOrder, error) {
	var order model.WithdrawalOrder
	err := r.db.WithContext(ctx).Where("kun_request_no = ?", requestNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
