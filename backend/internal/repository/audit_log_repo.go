package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
