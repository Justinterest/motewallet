package request

type UpdateSystemCurrencyConfigReq struct {
	CryptoCurrencies []string            `json:"crypto_currencies" binding:"required,min=1"`
	FiatCurrencies   []string            `json:"fiat_currencies" binding:"required,min=1"`
	CryptoChains     map[string][]string `json:"crypto_chains" binding:"required"`
	DefaultChains    map[string]string   `json:"default_chains"`
}
