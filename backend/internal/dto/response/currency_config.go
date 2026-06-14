package response

type SupportedCurrenciesResp struct {
	CryptoCurrencies []string `json:"crypto_currencies"`
	FiatCurrencies   []string `json:"fiat_currencies"`
}
