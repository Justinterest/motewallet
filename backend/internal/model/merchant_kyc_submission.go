package model

import "encoding/json"

type MerchantKycSubmission struct {
	BaseModel
	MerchantID   uint64          `gorm:"column:merchant_id;not null;index:idx_merchant_kyc_submissions_merchant_id" json:"merchant_id"`
	KunRequestNo string          `gorm:"column:kun_request_no;type:varchar(64);not null;uniqueIndex:uk_merchant_kyc_submissions_request_no" json:"kun_request_no"`
	KunAuthID    *string         `gorm:"column:kun_auth_id;type:varchar(128)" json:"kun_auth_id,omitempty"`
	Payload      json.RawMessage `gorm:"column:payload;type:json;not null" json:"payload"`
	Status       string          `gorm:"column:status;type:varchar(20);not null;default:SUBMITTED" json:"status"`
}

func (MerchantKycSubmission) TableName() string {
	return "merchant_kyc_submissions"
}
