package model

import "time"

type Merchant struct {
	SoftDeleteModel
	Email             string     `gorm:"column:email;type:varchar(128);not null;uniqueIndex:uk_merchants_email" json:"email"`
	PasswordHash      string     `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	KunSubCustomerNo  *string    `gorm:"column:kun_sub_customer_no;type:varchar(64);uniqueIndex:uk_merchants_kun_sub_customer_no" json:"kun_sub_customer_no,omitempty"`
	FeeTemplateID             *uint64    `gorm:"column:fee_template_id;index:idx_merchants_fee_template_id" json:"fee_template_id,omitempty"`
	SupportedCryptoCurrencies *string    `gorm:"column:supported_crypto_currencies;type:varchar(128)" json:"supported_crypto_currencies,omitempty"`
	SupportedFiatCurrencies   *string    `gorm:"column:supported_fiat_currencies;type:varchar(128)" json:"supported_fiat_currencies,omitempty"`
	Status                    string     `gorm:"column:status;type:varchar(20);not null;default:PENDING_AGREEMENT;index:idx_merchants_status" json:"status"`
	KycAuthID         *string    `gorm:"column:kyc_auth_id;type:varchar(128)" json:"kyc_auth_id,omitempty"`
	KycStatus         string     `gorm:"column:kyc_status;type:varchar(20);not null;default:NONE;index:idx_merchants_kyc_status" json:"kyc_status"`
	KycFailReason     *string    `gorm:"column:kyc_fail_reason;type:text" json:"kyc_fail_reason,omitempty"`
	KycSubmittedAt    *time.Time `gorm:"column:kyc_submitted_at;type:datetime(3)" json:"kyc_submitted_at,omitempty"`
	KycCompletedAt    *time.Time `gorm:"column:kyc_completed_at;type:datetime(3)" json:"kyc_completed_at,omitempty"`
	AgreementSignedAt *time.Time `gorm:"column:agreement_signed_at;type:datetime(3)" json:"agreement_signed_at,omitempty"`
	FrozenAt          *time.Time `gorm:"column:frozen_at;type:datetime(3)" json:"frozen_at,omitempty"`
}

func (Merchant) TableName() string {
	return "merchants"
}
