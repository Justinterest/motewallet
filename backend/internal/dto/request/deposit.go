package request

type GetDepositAddressReq struct {
	Currency string `form:"currency" binding:"required"`
	Chain    string `form:"chain" binding:"required"`
}
