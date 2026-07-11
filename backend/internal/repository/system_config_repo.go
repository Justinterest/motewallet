package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type SystemConfigRepository interface {
	GetByKey(ctx context.Context, key string) (string, error)
	Upsert(ctx context.Context, key, value, description string) error
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

func (r *systemConfigRepository) Upsert(ctx context.Context, key, value, description string) error {
	var config model.SystemConfig
	err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config = model.SystemConfig{
				ConfigKey:   key,
				ConfigValue: value,
				Description: &description,
			}
			return r.db.WithContext(ctx).Create(&config).Error
		}
		return err
	}
	updates := map[string]interface{}{"config_value": value}
	if description != "" {
		updates["description"] = description
	}
	return r.db.WithContext(ctx).Model(&config).Updates(updates).Error
}
