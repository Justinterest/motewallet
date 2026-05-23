package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet-withdrawal/backend/internal/model"
)

type WebhookLogRepository interface {
	Create(ctx context.Context, log *model.WebhookLog) error
	FindByEventIDAndTopic(ctx context.Context, eventID, eventTopic string) (*model.WebhookLog, error)
	UpdateProcessStatus(ctx context.Context, id uint64, status string, processError *string) error
}

type webhookLogRepository struct {
	db *gorm.DB
}

func NewWebhookLogRepository(db *gorm.DB) WebhookLogRepository {
	return &webhookLogRepository{db: db}
}

func (r *webhookLogRepository) Create(ctx context.Context, log *model.WebhookLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *webhookLogRepository) FindByEventIDAndTopic(ctx context.Context, eventID, eventTopic string) (*model.WebhookLog, error) {
	var log model.WebhookLog
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND event_topic = ?", eventID, eventTopic).
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *webhookLogRepository) UpdateProcessStatus(ctx context.Context, id uint64, status string, processError *string) error {
	updates := map[string]interface{}{
		"process_status": status,
		"process_error":  processError,
	}
	if status == "PROCESSED" || status == "FAILED" {
		updates["processed_at"] = gorm.Expr("NOW(3)")
	}
	return r.db.WithContext(ctx).Model(&model.WebhookLog{}).Where("id = ?", id).Updates(updates).Error
}
