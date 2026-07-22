package model

import "time"

type AdminUser struct {
	SoftDeleteModel
	Username          string     `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uk_admin_users_username" json:"username"`
	Email             string     `gorm:"column:email;type:varchar(128);not null;uniqueIndex:uk_admin_users_email" json:"email"`
	PasswordHash       string     `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	MustChangePassword bool       `gorm:"column:must_change_password;type:tinyint(1);not null;default:0" json:"must_change_password"`
	TotpSecret         *string    `gorm:"column:totp_secret;type:varchar(64)" json:"-"`
	TotpEnabled       bool       `gorm:"column:totp_enabled;type:tinyint(1);not null;default:0" json:"totp_enabled"`
	TotpPendingSecret *string    `gorm:"column:totp_pending_secret;type:varchar(64)" json:"-"`
	Role              string     `gorm:"column:role;type:varchar(20);not null;default:OPERATOR;index:idx_admin_users_role" json:"role"`
	Status            string     `gorm:"column:status;type:varchar(20);not null;default:ACTIVE" json:"status"`
	LastLoginAt       *time.Time `gorm:"column:last_login_at;type:datetime(3)" json:"last_login_at,omitempty"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}
