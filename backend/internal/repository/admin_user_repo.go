package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"motewallet-withdrawal/backend/internal/model"
)

type AdminUserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.AdminUser, error)
	FindByID(ctx context.Context, id uint64) (*model.AdminUser, error)
	UpdateLastLogin(ctx context.Context, id uint64) error
}

type adminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) FindByUsername(ctx context.Context, username string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *adminUserRepository) FindByID(ctx context.Context, id uint64) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *adminUserRepository) UpdateLastLogin(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.AdminUser{}).Where("id = ?", id).Update("last_login_at", now).Error
}
