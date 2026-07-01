package repository

import (
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"motewallet/internal/model"
)

type AdminWithdrawalListFilter struct {
	Page          int
	PageSize      int
	MerchantID    uint64
	Currency      string
	Status        string
	ReviewStatus  string
	Type          string
	MerchantEmail string
}

type AdminWithdrawalListRow struct {
	model.WithdrawalOrder
	MerchantEmail     string          `gorm:"column:merchant_email"`
	PlatformOrderID   string          `gorm:"column:platform_order_id"`
	Amount            decimal.Decimal `gorm:"column:amount"`
	PlatformFee       decimal.Decimal `gorm:"column:platform_fee"`
	Currency          string          `gorm:"column:currency"`
	TransactionStatus string          `gorm:"column:transaction_status"`
}

type WithdrawalOrderRepository interface {
	Create(ctx context.Context, order *model.WithdrawalOrder) error
	CreateWithDB(ctx context.Context, db *gorm.DB, order *model.WithdrawalOrder) error
	FindByID(ctx context.Context, id uint64) (*model.WithdrawalOrder, error)
	ListByMerchant(ctx context.Context, merchantID uint64, page, pageSize int) ([]*model.WithdrawalOrder, int64, error)
	ListPendingReview(ctx context.Context, page, pageSize int) ([]*model.WithdrawalOrder, int64, error)
	ListForAdmin(ctx context.Context, filter AdminWithdrawalListFilter) ([]AdminWithdrawalListRow, int64, error)
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

func (r *withdrawalOrderRepository) ListForAdmin(ctx context.Context, filter AdminWithdrawalListFilter) ([]AdminWithdrawalListRow, int64, error) {
	var rows []AdminWithdrawalListRow
	var total int64

	query := r.db.WithContext(ctx).
		Table("withdrawal_orders").
		Select(`
			withdrawal_orders.*,
			merchants.email AS merchant_email,
			transaction_records.platform_order_id AS platform_order_id,
			transaction_records.amount AS amount,
			transaction_records.platform_fee AS platform_fee,
			transaction_records.currency AS currency,
			transaction_records.status AS transaction_status
		`).
		Joins("JOIN merchants ON merchants.id = withdrawal_orders.merchant_id AND merchants.deleted_at IS NULL").
		Joins("JOIN transaction_records ON transaction_records.id = withdrawal_orders.transaction_record_id")

	if filter.MerchantID > 0 {
		query = query.Where("withdrawal_orders.merchant_id = ?", filter.MerchantID)
	}
	if filter.Currency != "" {
		query = query.Where("transaction_records.currency = ?", filter.Currency)
	}
	if filter.Status != "" {
		query = query.Where("transaction_records.status = ?", filter.Status)
	}
	if filter.ReviewStatus != "" {
		query = query.Where("withdrawal_orders.review_status = ?", filter.ReviewStatus)
	}
	if filter.Type != "" {
		query = query.Where("withdrawal_orders.withdrawal_type = ?", filter.Type)
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

	if err := query.Order("withdrawal_orders.id DESC").Offset(offset).Limit(pageSize).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
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
