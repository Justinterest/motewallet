package request

type SubmitCryptoWithdrawalReq struct {
	Currency  string `json:"currency" binding:"required"`
	Chain     string `json:"chain" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	ToAddress string `json:"to_address" binding:"required"`
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
