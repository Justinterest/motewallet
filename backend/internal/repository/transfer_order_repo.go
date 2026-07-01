package repository

import (
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"motewallet/internal/model"
)

type TransferOrderRepository interface {
	Create(ctx context.Context, order *model.TransferOrder) error
	FindByID(ctx context.Context, id uint64) (*model.TransferOrder, error)
	ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.TransferOrder, int64, error)
	FindByKunRequestNo(ctx context.Context, requestNo string) (*model.TransferOrder, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
	ListForAdmin(ctx context.Context, filter AdminTransferListFilter) ([]AdminTransferListRow, int64, error)
}

type AdminTransferListFilter struct {
	Page          int
	PageSize      int
	MerchantID    uint64
	Currency      string
	Status        string
	MerchantEmail string
}

type AdminTransferListRow struct {
	model.TransferOrder
	MerchantEmail     string          `gorm:"column:merchant_email"`
	PlatformOrderID   string          `gorm:"column:platform_order_id"`
	Amount            decimal.Decimal `gorm:"column:amount"`
	Currency          string          `gorm:"column:currency"`
	TransactionStatus string          `gorm:"column:transaction_status"`
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

func (r *transferOrderRepository) ListForAdmin(ctx context.Context, filter AdminTransferListFilter) ([]AdminTransferListRow, int64, error) {
	var rows []AdminTransferListRow
	var total int64

	query := r.db.WithContext(ctx).
		Table("transfer_orders").
		Select(`
			transfer_orders.*,
			merchants.email AS merchant_email,
			transaction_records.platform_order_id AS platform_order_id,
			transaction_records.amount AS amount,
			transaction_records.currency AS currency,
			transaction_records.status AS transaction_status
		`).
		Joins("JOIN merchants ON merchants.id = transfer_orders.merchant_id AND merchants.deleted_at IS NULL").
		Joins("JOIN transaction_records ON transaction_records.id = transfer_orders.transaction_record_id")

	if filter.MerchantID > 0 {
		query = query.Where("transfer_orders.merchant_id = ?", filter.MerchantID)
	}
	if filter.Currency != "" {
		query = query.Where("transaction_records.currency = ?", filter.Currency)
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
	if err := query.Order("transfer_orders.id DESC").Offset(offset).Limit(pageSize).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
