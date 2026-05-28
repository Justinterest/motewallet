package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type TransactionRecordRepository interface {
	Create(ctx context.Context, record *model.TransactionRecord) error
	FindByID(ctx context.Context, id uint64) (*model.TransactionRecord, error)
	FindByPlatformOrderID(ctx context.Context, platformOrderID string) (*model.TransactionRecord, error)
	UpdateStatus(ctx context.Context, id uint64, status string) error
	ListByMerchant(ctx context.Context, merchantID uint64, txType, status string, page, pageSize int) ([]*model.TransactionRecord, int64, error)
}

type transactionRecordRepository struct {
	db *gorm.DB
}

func NewTransactionRecordRepository(db *gorm.DB) TransactionRecordRepository {
	return &transactionRecordRepository{db: db}
}

func (r *transactionRecordRepository) Create(ctx context.Context, record *model.TransactionRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *transactionRecordRepository) FindByID(ctx context.Context, id uint64) (*model.TransactionRecord, error) {
	var record model.TransactionRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *transactionRecordRepository) FindByPlatformOrderID(ctx context.Context, platformOrderID string) (*model.TransactionRecord, error) {
	var record model.TransactionRecord
	err := r.db.WithContext(ctx).Where("platform_order_id = ?", platformOrderID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *transactionRecordRepository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	return r.db.WithContext(ctx).Model(&model.TransactionRecord{}).Where("id = ?", id).Update("status", status).Error
}

func (r *transactionRecordRepository) ListByMerchant(ctx context.Context, merchantID uint64, txType, status string, page, pageSize int) ([]*model.TransactionRecord, int64, error) {
	var records []*model.TransactionRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TransactionRecord{}).Where("merchant_id = ?", merchantID)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}
