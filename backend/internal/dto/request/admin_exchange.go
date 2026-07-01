package request

type AdminListExchangesReq struct {
	Page          int    `form:"page"`
	PageSize      int    `form:"page_size"`
	MerchantID    uint64 `form:"merchant_id"`
	Currency      string `form:"currency"`
	Status        string `form:"status"`
	MerchantEmail string `form:"merchant_email"`
}
