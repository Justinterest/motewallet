package repository

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"motewallet/internal/model"
)

type AdminDepositListFilter struct {
	Page          int
	PageSize      int
	MerchantID    uint64
	Currency      string
	Status        string
	MerchantEmail string
}

type AdminDepositListRow struct {
	model.DepositOrder
	MerchantEmail     string          `gorm:"column:merchant_email"`
	PlatformOrderID   string          `gorm:"column:platform_order_id"`
	Amount            decimal.Decimal `gorm:"column:amount"`
	TransactionStatus string          `gorm:"column:transaction_status"`
	TxnCompletedAt    *time.Time      `gorm:"column:txn_completed_at"`
}

type DepositOrderRepository interface {
	Create(ctx context.Context, order *model.DepositOrder) error
	CreateWithDB(ctx context.Context, db *gorm.DB, order *model.DepositOrder) error
	FindByID(ctx context.Context, id uint64) (*model.DepositOrder, error)
	ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.DepositOrder, int64, error)
	ListForAdmin(ctx context.Context, filter AdminDepositListFilter) ([]AdminDepositListRow, int64, error)
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

func (r *depositOrderRepository) ListForAdmin(ctx context.Context, filter AdminDepositListFilter) ([]AdminDepositListRow, int64, error) {
	var rows []AdminDepositListRow
	var total int64

	query := r.db.WithContext(ctx).
		Table("deposit_orders").
		Select(`
			deposit_orders.*,
			merchants.email AS merchant_email,
			transaction_records.platform_order_id AS platform_order_id,
			transaction_records.amount AS amount,
			transaction_records.status AS transaction_status,
			transaction_records.completed_at AS txn_completed_at
		`).
		Joins("JOIN merchants ON merchants.id = deposit_orders.merchant_id AND merchants.deleted_at IS NULL").
		Joins("JOIN transaction_records ON transaction_records.id = deposit_orders.transaction_record_id")

	if filter.MerchantID > 0 {
		query = query.Where("deposit_orders.merchant_id = ?", filter.MerchantID)
	}
	if filter.Currency != "" {
		query = query.Where("deposit_orders.currency = ?", filter.Currency)
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

	if err := query.Order("deposit_orders.id DESC").Offset(offset).Limit(pageSize).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *depositOrderRepository) FindByKunOrderID(ctx context.Context, kunOrderID string) (*model.DepositOrder, error) {
	var order model.DepositOrder
	err := r.db.WithContext(ctx).Where("kun_order_id = ?", kunOrderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
