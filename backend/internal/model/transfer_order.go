package model

import "time"

type TransferOrder struct {
	BaseModel
	TransactionRecordID uint64     `gorm:"column:transaction_record_id;not null;uniqueIndex:uk_transfer_orders_txn_record" json:"transaction_record_id"`
	MerchantID          uint64     `gorm:"column:merchant_id;not null;index:idx_transfer_orders_merchant_id" json:"merchant_id"`
	FromAccountType     string     `gorm:"column:from_account_type;type:varchar(10);not null" json:"from_account_type"`
	ToAccountType       string     `gorm:"column:to_account_type;type:varchar(10);not null" json:"to_account_type"`
	KunOrderID          *string    `gorm:"column:kun_order_id;type:varchar(128)" json:"kun_order_id,omitempty"`
	KunRequestNo        *string    `gorm:"column:kun_request_no;type:varchar(64);uniqueIndex:uk_transfer_orders_kun_request_no" json:"kun_request_no,omitempty"`
	CompletedAt         *time.Time `gorm:"column:completed_at;type:datetime(3)" json:"completed_at,omitempty"`
}

func (TransferOrder) TableName() string {
	return "transfer_orders"
}
