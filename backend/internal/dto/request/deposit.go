package request

type GetDepositAddressReq struct {
	Currency string `form:"currency" binding:"required"`
	Chain    string `form:"chain" binding:"required"`
}

type ListDepositOrdersReq struct {
	Currency string `form:"currency"`
	Chain    string `form:"chain"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
