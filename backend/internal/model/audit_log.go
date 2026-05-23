package model

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt    time.Time       `gorm:"type:datetime(3);not null;autoCreateTime;index:idx_audit_logs_created_at" json:"created_at"`
	OperatorID   uint64          `gorm:"column:operator_id;not null;index:idx_audit_logs_operator" json:"operator_id"`
	OperatorType string          `gorm:"column:operator_type;type:varchar(10);not null;index:idx_audit_logs_operator" json:"operator_type"`
	Action       string          `gorm:"column:action;type:varchar(64);not null;index:idx_audit_logs_action" json:"action"`
	TargetType   *string         `gorm:"column:target_type;type:varchar(32);index:idx_audit_logs_target" json:"target_type,omitempty"`
	TargetID     *string         `gorm:"column:target_id;type:varchar(64);index:idx_audit_logs_target" json:"target_id,omitempty"`
	Detail       json.RawMessage `gorm:"column:detail;type:json" json:"detail,omitempty"`
	IPAddress    *string         `gorm:"column:ip_address;type:varchar(45)" json:"ip_address,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
