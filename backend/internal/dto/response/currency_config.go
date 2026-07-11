package response

type SupportedCurrenciesResp struct {
	CryptoCurrencies []string            `json:"crypto_currencies"`
	FiatCurrencies   []string            `json:"fiat_currencies"`
	CryptoChains     map[string][]string `json:"crypto_chains"`
	DefaultChains    map[string]string   `json:"default_chains"`
}

type SystemCurrencyConfigResp struct {
	CryptoCurrencies []string            `json:"crypto_currencies"`
	FiatCurrencies   []string            `json:"fiat_currencies"`
	CryptoChains     map[string][]string `json:"crypto_chains"`
	DefaultChains    map[string]string   `json:"default_chains"`
	CatalogChains    map[string][]string `json:"catalog_chains"`
	AllCrypto        []string            `json:"all_crypto"`
	AllFiat          []string            `json:"all_fiat"`
}
