package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type DepositOrderRepository interface {
	Create(ctx context.Context, order *model.DepositOrder) error
	CreateWithDB(ctx context.Context, db *gorm.DB, order *model.DepositOrder) error
	FindByID(ctx context.Context, id uint64) (*model.DepositOrder, error)
	ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.DepositOrder, int64, error)
	FindByKunOrderID(ctx context.Context, kunOrderID string) (*model.DepositOrder, error)
}

type depositOrderRepository struct {
	db *gorm.DB
}

func NewDepositOrderRepository(db *gorm.DB) DepositOrderRepository {
	return &depositOrderRepository{db: db}
}

func (r *depositOrderRepository) Create(ctx context.Context, order *model.DepositOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *depositOrderRepository) CreateWithDB(ctx context.Context, db *gorm.DB, order *model.DepositOrder) error {
	return db.WithContext(ctx).Create(order).Error
}

func (r *depositOrderRepository) FindByID(ctx context.Context, id uint64) (*model.DepositOrder, error) {
	var order model.DepositOrder
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *depositOrderRepository) ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.DepositOrder, int64, error) {
	var orders []*model.DepositOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&model.DepositOrder{}).Where("merchant_id = ?", merchantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *depositOrderRepository) FindByKunOrderID(ctx context.Context, kunOrderID string) (*model.DepositOrder, error) {
	var order model.DepositOrder
	err := r.db.WithContext(ctx).Where("kun_order_id = ?", kunOrderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
