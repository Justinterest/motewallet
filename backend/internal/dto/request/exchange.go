package request

type ExchangePreviewReq struct {
	FromCurrency string `json:"from_currency" binding:"required"`
	ToCurrency   string `json:"to_currency" binding:"required"`
	FromAmount   string `json:"from_amount" binding:"required"`
}

type CreateExchangeOrderReq struct {
	FromCurrency string `json:"from_currency" binding:"required"`
	ToCurrency   string `json:"to_currency" binding:"required"`
	FromAmount   string `json:"from_amount" binding:"required"`
}
