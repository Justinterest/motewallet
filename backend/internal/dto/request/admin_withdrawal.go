package request

type AdminListWithdrawalsReq struct {
	Page          int    `form:"page"`
	PageSize      int    `form:"page_size"`
	MerchantID    uint64 `form:"merchant_id"`
	Currency      string `form:"currency"`
	Status        string `form:"status"`
	ReviewStatus  string `form:"review_status"`
	Type          string `form:"type"`
	MerchantEmail string `form:"merchant_email"`
}
