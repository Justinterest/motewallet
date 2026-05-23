package request

type TransferReq struct {
	FromAccountType string `json:"from_account_type" binding:"required"`
	ToAccountType   string `json:"to_account_type" binding:"required"`
	Currency        string `json:"currency" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
}
