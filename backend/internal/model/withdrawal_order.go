package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type WithdrawalOrder struct {
	BaseModel
	TransactionRecordID uint64           `gorm:"column:transaction_record_id;not null;uniqueIndex:uk_withdrawal_orders_txn_record" json:"transaction_record_id"`
	MerchantID          uint64           `gorm:"column:merchant_id;not null;index:idx_withdrawal_orders_merchant_id" json:"merchant_id"`
	WithdrawalType      string           `gorm:"column:withdrawal_type;type:varchar(20);not null;index:idx_withdrawal_orders_withdrawal_type" json:"withdrawal_type"`
	CryptoAddressID     *uint64          `gorm:"column:crypto_address_id" json:"crypto_address_id,omitempty"`
	BankAccountID       *uint64          `gorm:"column:bank_account_id" json:"bank_account_id,omitempty"`
	ToAddress           *string          `gorm:"column:to_address;type:varchar(255)" json:"to_address,omitempty"`
	Chain               *string          `gorm:"column:chain;type:varchar(20)" json:"chain,omitempty"`
	TxID                *string          `gorm:"column:tx_id;type:varchar(128)" json:"tx_id,omitempty"`
	TransferType        *string          `gorm:"column:transfer_type;type:varchar(10)" json:"transfer_type,omitempty"`
	Purpose             *string          `gorm:"column:purpose;type:varchar(10)" json:"purpose,omitempty"`
	ReviewStatus        string           `gorm:"column:review_status;type:varchar(20);not null;default:PENDING_REVIEW;index:idx_withdrawal_orders_review_status" json:"review_status"`
	ReviewerID          *uint64          `gorm:"column:reviewer_id;index:idx_withdrawal_orders_reviewer_id" json:"reviewer_id,omitempty"`
	ReviewerType        *string          `gorm:"column:reviewer_type;type:varchar(10)" json:"reviewer_type,omitempty"`
	ReviewedAt          *time.Time       `gorm:"column:reviewed_at;type:datetime(3)" json:"reviewed_at,omitempty"`
	ReviewRemark        *string          `gorm:"column:review_remark;type:varchar(500)" json:"review_remark,omitempty"`
	KunOrderID          *string          `gorm:"column:kun_order_id;type:varchar(128)" json:"kun_order_id,omitempty"`
	KunRequestNo        *string          `gorm:"column:kun_request_no;type:varchar(64);uniqueIndex:uk_withdrawal_orders_kun_request_no" json:"kun_request_no,omitempty"`
	KunFee              *decimal.Decimal `gorm:"column:kun_fee;type:decimal(28,8);default:0" json:"kun_fee,omitempty"`
	KunFeeCurrency      *string          `gorm:"column:kun_fee_currency;type:varchar(10)" json:"kun_fee_currency,omitempty"`
	KunSubmittedAt      *time.Time       `gorm:"column:kun_submitted_at;type:datetime(3)" json:"kun_submitted_at,omitempty"`
	FailedReason        *string          `gorm:"column:failed_reason;type:varchar(500)" json:"failed_reason,omitempty"`
	CompletedAt         *time.Time       `gorm:"column:completed_at;type:datetime(3)" json:"completed_at,omitempty"`
}

func (WithdrawalOrder) TableName() string {
	return "withdrawal_orders"
}
