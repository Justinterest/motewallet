package repository

import (
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"motewallet/internal/model"
)

type ExchangeOrderRepository interface {
	Create(ctx context.Context, order *model.ExchangeOrder) error
	FindByID(ctx context.Context, id uint64) (*model.ExchangeOrder, error)
	ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.ExchangeOrder, int64, error)
	FindByKunRequestNo(ctx context.Context, requestNo string) (*model.ExchangeOrder, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
	ListForAdmin(ctx context.Context, filter AdminExchangeListFilter) ([]AdminExchangeListRow, int64, error)
}

type AdminExchangeListFilter struct {
	Page          int
	PageSize      int
	MerchantID    uint64
	Currency      string
	Status        string
	MerchantEmail string
}

type AdminExchangeListRow struct {
	model.ExchangeOrder
	MerchantEmail     string          `gorm:"column:merchant_email"`
	PlatformOrderID   string          `gorm:"column:platform_order_id"`
	PlatformFee       decimal.Decimal `gorm:"column:platform_fee"`
	TransactionStatus string          `gorm:"column:transaction_status"`
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

func (r *exchangeOrderRepository) ListForAdmin(ctx context.Context, filter AdminExchangeListFilter) ([]AdminExchangeListRow, int64, error) {
	var rows []AdminExchangeListRow
	var total int64

	query := r.db.WithContext(ctx).
		Table("exchange_orders").
		Select(`
			exchange_orders.*,
			merchants.email AS merchant_email,
			transaction_records.platform_order_id AS platform_order_id,
			transaction_records.platform_fee AS platform_fee,
			transaction_records.status AS transaction_status
		`).
		Joins("JOIN merchants ON merchants.id = exchange_orders.merchant_id AND merchants.deleted_at IS NULL").
		Joins("JOIN transaction_records ON transaction_records.id = exchange_orders.transaction_record_id")

	if filter.MerchantID > 0 {
		query = query.Where("exchange_orders.merchant_id = ?", filter.MerchantID)
	}
	if filter.Currency != "" {
		query = query.Where("(exchange_orders.from_currency = ? OR exchange_orders.to_currency = ?)", filter.Currency, filter.Currency)
	}
	if filter.Status != "" {
		query = query.Where("transaction_records.status = ?", filter.Status)
	}
	if filter.MerchantEmail != "" {
		query = query.Where("merchants.email LIKE ?", "%"+filter.MerchantEmail+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Order("exchange_orders.id DESC").Offset(offset).Limit(pageSize).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
