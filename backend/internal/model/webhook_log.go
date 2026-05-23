package model

import (
	"encoding/json"
	"time"
)

type WebhookLog struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt     time.Time       `gorm:"type:datetime(3);not null;autoCreateTime" json:"created_at"`
	EventID       string          `gorm:"column:event_id;type:varchar(128);not null;uniqueIndex:uk_webhook_logs_event" json:"event_id"`
	EventTopic    string          `gorm:"column:event_topic;type:varchar(64);not null;uniqueIndex:uk_webhook_logs_event" json:"event_topic"`
	EventType     string          `gorm:"column:event_type;type:varchar(20);not null" json:"event_type"`
	CustomerNo    *string         `gorm:"column:customer_no;type:varchar(64);index:idx_webhook_logs_customer_no" json:"customer_no,omitempty"`
	RawData       json.RawMessage `gorm:"column:raw_data;type:json;not null" json:"raw_data"`
	ProcessStatus string          `gorm:"column:process_status;type:varchar(20);not null;default:PENDING;index:idx_webhook_logs_process_status" json:"process_status"`
	ProcessError  *string         `gorm:"column:process_error;type:text" json:"process_error,omitempty"`
	ProcessedAt   *time.Time      `gorm:"column:processed_at;type:datetime(3)" json:"processed_at,omitempty"`
	ReceivedAt    time.Time       `gorm:"column:received_at;type:datetime(3);not null;index:idx_webhook_logs_received_at" json:"received_at"`
}

func (WebhookLog) TableName() string {
	return "webhook_logs"
}
