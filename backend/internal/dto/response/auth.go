package response

import "time"

const (
	AuthStatusSuccess          = "SUCCESS"
	AuthStatusRequires2FA      = "REQUIRES_2FA"
	AuthStatusRequires2FASetup = "REQUIRES_2FA_SETUP"
)

type MerchantInfoResp struct {
	ID            uint64    `json:"id"`
	Email         string    `json:"email"`
	Status        string    `json:"status"`
	KycStatus     string    `json:"kyc_status"`
	TotpEnabled   bool      `json:"totp_enabled"`
	FeeTemplateID *uint64   `json:"fee_template_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// AuthChallengeResp is returned by login/register when further 2FA steps are required,
// or when authentication completes successfully.
type AuthChallengeResp struct {
	Status     string            `json:"status"`
	TempToken  string            `json:"temp_token,omitempty"`
	TotpSecret string            `json:"totp_secret,omitempty"`
	TotpURI    string            `json:"totp_uri,omitempty"`
	Merchant   *MerchantInfoResp `json:"merchant,omitempty"`
}

type TotpStatusResp struct {
	Enabled bool `json:"enabled"`
}

type TotpSetupResp struct {
	TotpSecret string `json:"totp_secret"`
	TotpURI    string `json:"totp_uri"`
}
