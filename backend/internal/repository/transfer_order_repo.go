package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet-withdrawal/backend/internal/model"
)

type TransferOrderRepository interface {
	Create(ctx context.Context, order *model.TransferOrder) error
	FindByID(ctx context.Context, id uint64) (*model.TransferOrder, error)
	ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.TransferOrder, int64, error)
	FindByKunRequestNo(ctx context.Context, requestNo string) (*model.TransferOrder, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
}

type transferOrderRepository struct {
	db *gorm.DB
}

func NewTransferOrderRepository(db *gorm.DB) TransferOrderRepository {
	return &transferOrderRepository{db: db}
}

func (r *transferOrderRepository) Create(ctx context.Context, order *model.TransferOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *transferOrderRepository) FindByID(ctx context.Context, id uint64) (*model.TransferOrder, error) {
	var order model.TransferOrder
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *transferOrderRepository) ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.TransferOrder, int64, error) {
	var orders []*model.TransferOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TransferOrder{}).Where("merchant_id = ?", merchantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *transferOrderRepository) FindByKunRequestNo(ctx context.Context, requestNo string) (*model.TransferOrder, error) {
	var order model.TransferOrder
	err := r.db.WithContext(ctx).Where("kun_request_no = ?", requestNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *transferOrderRepository) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.TransferOrder{}).Where("id = ?", id).Updates(fields).Error
}
