package request

type SubmitCryptoWithdrawalReq struct {
	Currency        string `json:"currency" binding:"required"`
	CryptoAddressID uint64 `json:"crypto_address_id" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
}

type SubmitFiatWithdrawalReq struct {
	Currency      string `json:"currency" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
	BankAccountID uint64 `json:"bank_account_id" binding:"required"`
	Purpose       string `json:"purpose" binding:"required"`
	Postscript    string `json:"postscript" binding:"required"`
}

type ReviewWithdrawalReq struct {
	Reason string `json:"reason"`
}

type WithdrawalFeePreviewReq struct {
	Type            string `json:"type" binding:"required,oneof=CRYPTO FIAT"`
	Currency        string `json:"currency" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	CryptoAddressID uint64 `json:"crypto_address_id"`
	BankAccountID   uint64 `json:"bank_account_id"`
}
