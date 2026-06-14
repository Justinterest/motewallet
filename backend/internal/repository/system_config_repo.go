package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type SystemConfigRepository interface {
	GetByKey(ctx context.Context, key string) (string, error)
}

type systemConfigRepository struct {
	db *gorm.DB
}

func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepository {
	return &systemConfigRepository{db: db}
}

func (r *systemConfigRepository) GetByKey(ctx context.Context, key string) (string, error) {
	var config model.SystemConfig
	err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return config.ConfigValue, nil
}
