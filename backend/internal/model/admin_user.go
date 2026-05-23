package model

import "time"

type AdminUser struct {
	SoftDeleteModel
	Username     string     `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uk_admin_users_username" json:"username"`
	Email        string     `gorm:"column:email;type:varchar(128);not null;uniqueIndex:uk_admin_users_email" json:"email"`
	PasswordHash string     `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Role         string     `gorm:"column:role;type:varchar(20);not null;default:OPERATOR;index:idx_admin_users_role" json:"role"`
	Status       string     `gorm:"column:status;type:varchar(20);not null;default:ACTIVE" json:"status"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at;type:datetime(3)" json:"last_login_at,omitempty"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}
