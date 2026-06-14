package response

import "time"

type AdminMerchantListItem struct {
	ID                uint64     `json:"id"`
	Email             string     `json:"email"`
	Status            string     `json:"status"`
	KycStatus         string     `json:"kyc_status"`
	FeeTemplateID     *uint64    `json:"fee_template_id,omitempty"`
	KunSubCustomerNo  *string    `json:"kun_sub_customer_no,omitempty"`
	AgreementSignedAt *time.Time `json:"agreement_signed_at,omitempty"`
	KycSubmittedAt    *time.Time `json:"kyc_submitted_at,omitempty"`
	KycCompletedAt    *time.Time `json:"kyc_completed_at,omitempty"`
	FrozenAt          *time.Time `json:"frozen_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type AdminMerchantDetailResp struct {
	ID                uint64              `json:"id"`
	Email             string              `json:"email"`
	Status            string              `json:"status"`
	KycStatus         string              `json:"kyc_status"`
	KycFailReason     *string             `json:"kyc_fail_reason,omitempty"`
	FeeTemplateID     *uint64             `json:"fee_template_id,omitempty"`
	FeeTemplateName   *string             `json:"fee_template_name,omitempty"`
	KunSubCustomerNo  *string             `json:"kun_sub_customer_no,omitempty"`
	AgreementSignedAt *time.Time          `json:"agreement_signed_at,omitempty"`
	KycSubmittedAt    *time.Time          `json:"kyc_submitted_at,omitempty"`
	KycCompletedAt    *time.Time          `json:"kyc_completed_at,omitempty"`
	FrozenAt          *time.Time          `json:"frozen_at,omitempty"`
	CreatedAt                 time.Time           `json:"created_at"`
	Wallets                   []WalletBalanceResp `json:"wallets"`
	SupportedCryptoCurrencies []string            `json:"supported_crypto_currencies"`
	SupportedFiatCurrencies   []string            `json:"supported_fiat_currencies"`
	AvailableCryptoCurrencies []string            `json:"available_crypto_currencies"`
	AvailableFiatCurrencies   []string            `json:"available_fiat_currencies"`
}
