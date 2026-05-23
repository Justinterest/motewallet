package response

import "time"

type MerchantInfoResp struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	KycStatus string    `json:"kyc_status"`
	CreatedAt time.Time `json:"created_at"`
}
