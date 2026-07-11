package request

type AdminListMerchantsReq struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status    string `form:"status"`
	KycStatus string `form:"kyc_status"`
	Search    string `form:"search"`
}

type UpdateMerchantStatusReq struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE FROZEN"`
}

type UpdateMerchantFeeTemplateReq struct {
	FeeTemplateID uint64 `json:"fee_template_id" binding:"required"`
}

type KycRejectReq struct {
	Reason string `json:"reason" binding:"required"`
}

type UpdateMerchantSupportedCurrenciesReq struct {
	CryptoCurrencies []string            `json:"crypto_currencies" binding:"required,min=1"`
	FiatCurrencies   []string            `json:"fiat_currencies" binding:"required,min=1"`
	CryptoChains     map[string][]string `json:"crypto_chains"`
	DefaultChains    map[string]string   `json:"default_chains"`
}
