package response

type CryptoAddressResp struct {
	ID       uint64 `json:"id"`
	Currency string `json:"currency"`
	Chain    string `json:"chain"`
	Address  string `json:"address"`
	Alias    string `json:"alias"`
	Status   string `json:"status"`
}

type BankAccountResp struct {
	ID           uint64 `json:"id"`
	Currency     string `json:"currency"`
	BankName     string `json:"bank_name"`
	BankCountry  string `json:"bank_country"`
	SwiftCode    string `json:"swift_code"`
	AccountName  string `json:"account_name"`
	TransferType string `json:"transfer_type"`
	Status       string `json:"status"`
}
