package response

import "time"

type AgreementResp struct {
	ID         string `json:"id"`
	ProtocolID string `json:"protocol_id"`
	Title      string `json:"title"`
	Version    string `json:"version,omitempty"`
	URL        string `json:"url"`
	Required   bool   `json:"required"`
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
