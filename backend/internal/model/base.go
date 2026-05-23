package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime(3);not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3);not null;autoUpdateTime" json:"updated_at"`
}

type SoftDeleteModel struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"type:datetime(3);index" json:"-"`
}
