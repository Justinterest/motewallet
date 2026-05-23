package request

type AddCryptoAddressReq struct {
	Currency string `json:"currency" binding:"required"`
	Chain    string `json:"chain" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Alias    string `json:"alias"`
}

type AddBankAccountReq struct {
	Currency     string `json:"currency" binding:"required"`
	BankName     string `json:"bank_name" binding:"required"`
	BankCountry  string `json:"bank_country" binding:"required"`
	SwiftCode    string `json:"swift_code" binding:"required"`
	AccountName  string `json:"account_name" binding:"required"`
	AccountNo    string `json:"account_no" binding:"required"`
	TransferType string `json:"transfer_type" binding:"required"`
}
