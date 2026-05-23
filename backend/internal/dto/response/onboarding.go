package response

import "time"

type AgreementResp struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Required bool   `json:"required"`
}

type AgreementListResp struct {
	Agreements []AgreementResp `json:"agreements"`
	Signed     bool            `json:"signed"`
}

type KycStatusResp struct {
	Status         string     `json:"status"`
	KycStatus      string     `json:"kyc_status"`
	KycFailReason  *string    `json:"kyc_fail_reason,omitempty"`
	KycSubmittedAt *time.Time `json:"kyc_submitted_at,omitempty"`
	KycCompletedAt *time.Time `json:"kyc_completed_at,omitempty"`
	AgreementSignedAt *time.Time `json:"agreement_signed_at,omitempty"`
}
