package request

type AddCryptoAddressReq struct {
	Currency string `json:"currency" binding:"required"`
	Chain    string `json:"chain" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Alias    string `json:"alias" binding:"required"`
}

type AddBankAccountReq struct {
	Currency           string `json:"currency" binding:"required"`
	BankName           string `json:"bank_name" binding:"required"`
	BankCountry        string `json:"bank_country" binding:"required"`
	SwiftCode          string `json:"swift_code"`
	AccountName        string `json:"account_name" binding:"required"`
	AccountNo          string `json:"account_no" binding:"required"`
	TransferType       string `json:"transfer_type" binding:"required"`
	AccountType        string `json:"account_type"`
	PayeeCountryCode   string `json:"payee_country_code"`
	PayeeAddress       string `json:"payee_address"`
	PayeeAddressSecond string `json:"payee_address_second"`
	BankCode           string `json:"bank_code"`
	BankAddress        string `json:"bank_address"`
	MiddleSwiftCode    string `json:"middle_swift_code"`
}
