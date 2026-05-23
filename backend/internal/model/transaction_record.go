package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionRecord struct {
	BaseModel
	PlatformOrderID     string           `gorm:"column:platform_order_id;type:varchar(64);not null;uniqueIndex:uk_transaction_records_platform_order_id" json:"platform_order_id"`
	MerchantID          uint64           `gorm:"column:merchant_id;not null;index:idx_transaction_records_merchant_id;index:idx_transaction_records_merchant_type;index:idx_transaction_records_merchant_status" json:"merchant_id"`
	Type                string           `gorm:"column:type;type:varchar(20);not null;index:idx_transaction_records_type;index:idx_transaction_records_merchant_type" json:"type"`
	SubType             *string          `gorm:"column:sub_type;type:varchar(30)" json:"sub_type,omitempty"`
	Amount              decimal.Decimal  `gorm:"column:amount;type:decimal(28,8);not null" json:"amount"`
	Currency            string           `gorm:"column:currency;type:varchar(10);not null" json:"currency"`
	PlatformFee         decimal.Decimal  `gorm:"column:platform_fee;type:decimal(28,8);not null;default:0" json:"platform_fee"`
	PlatformFeeCurrency *string          `gorm:"column:platform_fee_currency;type:varchar(10)" json:"platform_fee_currency,omitempty"`
	ActualAmount        *decimal.Decimal `gorm:"column:actual_amount;type:decimal(28,8)" json:"actual_amount,omitempty"`
	Remark              *string          `gorm:"column:remark;type:varchar(500)" json:"remark,omitempty"`
	Status              string           `gorm:"column:status;type:varchar(20);not null;default:PENDING;index:idx_transaction_records_status;index:idx_transaction_records_merchant_status" json:"status"`
	CompletedAt         *time.Time       `gorm:"column:completed_at;type:datetime(3)" json:"completed_at,omitempty"`
	FailedReason        *string          `gorm:"column:failed_reason;type:varchar(500)" json:"failed_reason,omitempty"`
}

func (TransactionRecord) TableName() string {
	return "transaction_records"
}
