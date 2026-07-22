package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type AdminUserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.AdminUser, error)
	FindByEmail(ctx context.Context, email string) (*model.AdminUser, error)
	FindByID(ctx context.Context, id uint64) (*model.AdminUser, error)
	List(ctx context.Context) ([]model.AdminUser, error)
	Create(ctx context.Context, user *model.AdminUser) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error
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

func (r *adminUserRepository) FindByEmail(ctx context.Context, email string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
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

func (r *adminUserRepository) List(ctx context.Context) ([]model.AdminUser, error) {
	var users []model.AdminUser
	err := r.db.WithContext(ctx).Order("id ASC").Find(&users).Error
	return users, err
}

func (r *adminUserRepository) Create(ctx context.Context, user *model.AdminUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *adminUserRepository) UpdateFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.AdminUser{}).Where("id = ?", id).Updates(fields).Error
}

func (r *adminUserRepository) UpdateLastLogin(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.AdminUser{}).Where("id = ?", id).Update("last_login_at", now).Error
}
