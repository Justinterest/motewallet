package model

import "time"

type SystemAnnouncement struct {
	SoftDeleteModel
	Title       string     `gorm:"column:title;type:varchar(128);not null" json:"title"`
	Content     string     `gorm:"column:content;type:text;not null" json:"content"`
	Status      string     `gorm:"column:status;type:varchar(20);not null;default:DRAFT;index:idx_system_announcements_status" json:"status"`
	PublishedAt *time.Time `gorm:"column:published_at;type:datetime(3);index:idx_system_announcements_published_at" json:"published_at,omitempty"`
	CreatedBy   uint64     `gorm:"column:created_by;not null" json:"created_by"`
}

func (SystemAnnouncement) TableName() string {
	return "system_announcements"
}
