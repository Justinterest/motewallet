package model

import "time"

type DepositOrder struct {
	BaseModel
	TransactionRecordID uint64     `gorm:"column:transaction_record_id;not null;uniqueIndex:uk_deposit_orders_txn_record" json:"transaction_record_id"`
	MerchantID          uint64     `gorm:"column:merchant_id;not null;index:idx_deposit_orders_merchant_id" json:"merchant_id"`
	Currency            string     `gorm:"column:currency;type:varchar(10);not null" json:"currency"`
	Chain               string     `gorm:"column:chain;type:varchar(20);not null" json:"chain"`
	ToAddress           string     `gorm:"column:to_address;type:varchar(255);not null" json:"to_address"`
	FromAddress         *string    `gorm:"column:from_address;type:varchar(255)" json:"from_address,omitempty"`
	TxID                *string    `gorm:"column:tx_id;type:varchar(128);index:idx_deposit_orders_tx_id" json:"tx_id,omitempty"`
	KunOrderID          *string    `gorm:"column:kun_order_id;type:varchar(128);index:idx_deposit_orders_kun_order_id" json:"kun_order_id,omitempty"`
	Confirmations       *uint      `gorm:"column:confirmations;type:int unsigned" json:"confirmations,omitempty"`
	CompletedAt         *time.Time `gorm:"column:completed_at;type:datetime(3)" json:"completed_at,omitempty"`
}

func (DepositOrder) TableName() string {
	return "deposit_orders"
}
