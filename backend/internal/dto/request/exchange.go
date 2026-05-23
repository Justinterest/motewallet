package request

type GetExchangeQuoteReq struct {
	FromCurrency string `json:"from_currency" binding:"required"`
	ToCurrency   string `json:"to_currency" binding:"required"`
	FromAmount   string `json:"from_amount" binding:"required"`
}

type CreateExchangeOrderReq struct {
	QuoteId      string `json:"quote_id" binding:"required"`
	FromCurrency string `json:"from_currency" binding:"required"`
	ToCurrency   string `json:"to_currency" binding:"required"`
	FromAmount   string `json:"from_amount" binding:"required"`
}

type Create1to1OrderReq struct {
	FromCurrency string `json:"from_currency" binding:"required"`
	ToCurrency   string `json:"to_currency" binding:"required"`
	FromAmount   string `json:"from_amount" binding:"required"`
}
