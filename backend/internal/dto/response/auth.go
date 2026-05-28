package response

import "time"

type MerchantInfoResp struct {
	ID            uint64    `json:"id"`
	Email         string    `json:"email"`
	Status        string    `json:"status"`
	KycStatus     string    `json:"kyc_status"`
	FeeTemplateID *uint64   `json:"fee_template_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
