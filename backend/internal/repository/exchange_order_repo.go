package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type ExchangeOrderRepository interface {
	Create(ctx context.Context, order *model.ExchangeOrder) error
	FindByID(ctx context.Context, id uint64) (*model.ExchangeOrder, error)
	ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.ExchangeOrder, int64, error)
	FindByKunRequestNo(ctx context.Context, requestNo string) (*model.ExchangeOrder, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
}

type exchangeOrderRepository struct {
	db *gorm.DB
}

func NewExchangeOrderRepository(db *gorm.DB) ExchangeOrderRepository {
	return &exchangeOrderRepository{db: db}
}

func (r *exchangeOrderRepository) Create(ctx context.Context, order *model.ExchangeOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *exchangeOrderRepository) FindByID(ctx context.Context, id uint64) (*model.ExchangeOrder, error) {
	var order model.ExchangeOrder
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *exchangeOrderRepository) ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.ExchangeOrder, int64, error) {
	var orders []*model.ExchangeOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ExchangeOrder{}).Where("merchant_id = ?", merchantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *exchangeOrderRepository) FindByKunRequestNo(ctx context.Context, requestNo string) (*model.ExchangeOrder, error) {
	var order model.ExchangeOrder
	err := r.db.WithContext(ctx).Where("kun_request_no = ?", requestNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *exchangeOrderRepository) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.ExchangeOrder{}).Where("id = ?", id).Updates(fields).Error
}
