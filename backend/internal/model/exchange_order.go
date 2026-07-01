package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeOrder struct {
	BaseModel
	TransactionRecordID uint64           `gorm:"column:transaction_record_id;not null;uniqueIndex:uk_exchange_orders_txn_record" json:"transaction_record_id"`
	MerchantID          uint64           `gorm:"column:merchant_id;not null;index:idx_exchange_orders_merchant_id" json:"merchant_id"`
	ExchangeType        string           `gorm:"column:exchange_type;type:varchar(20);not null;index:idx_exchange_orders_exchange_type" json:"exchange_type"`
	FromCurrency        string           `gorm:"column:from_currency;type:varchar(10);not null" json:"from_currency"`
	FromAmount          decimal.Decimal  `gorm:"column:from_amount;type:decimal(28,8);not null" json:"from_amount"`
	ToCurrency          string           `gorm:"column:to_currency;type:varchar(10);not null" json:"to_currency"`
	ToAmount            *decimal.Decimal `gorm:"column:to_amount;type:decimal(28,8)" json:"to_amount,omitempty"`
	ExchangeRate        *decimal.Decimal `gorm:"column:exchange_rate;type:decimal(20,10)" json:"exchange_rate,omitempty"`
	QuoteID             *string          `gorm:"column:quote_id;type:varchar(128)" json:"quote_id,omitempty"`
	AutoTransfer        *string          `gorm:"column:auto_transfer;type:varchar(3);default:NO" json:"auto_transfer,omitempty"`
	KunOrderID          *string          `gorm:"column:kun_order_id;type:varchar(128);index:idx_exchange_orders_kun_order_id" json:"kun_order_id,omitempty"`
	KunRequestNo        *string          `gorm:"column:kun_request_no;type:varchar(64);uniqueIndex:uk_exchange_orders_kun_request_no" json:"kun_request_no,omitempty"`
	KunFee              *decimal.Decimal `gorm:"column:kun_fee;type:decimal(28,8);default:0" json:"kun_fee,omitempty"`
	KunFeeCurrency      *string          `gorm:"column:kun_fee_currency;type:varchar(10)" json:"kun_fee_currency,omitempty"`
	CompletedAt         *time.Time       `gorm:"column:completed_at;type:datetime(3)" json:"completed_at,omitempty"`
	FailReason          *string          `gorm:"column:fail_reason;type:text" json:"fail_reason,omitempty"`
}

func (ExchangeOrder) TableName() string {
	return "exchange_orders"
}
