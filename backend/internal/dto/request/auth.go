package request

type RegisterReq struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=8,max=64"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
}

type SendVerificationCodeReq struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TotpVerifyReq struct {
	TempToken string `json:"temp_token" binding:"required"`
	Code      string `json:"code" binding:"required,len=6"`
}

type TotpSetupConfirmReq struct {
	TempToken string `json:"temp_token" binding:"required"`
	Code      string `json:"code" binding:"required,len=6"`
}

type TotpRebindPrepareReq struct {
	CurrentCode string `json:"current_code" binding:"required,len=6"`
}

type TotpRebindConfirmReq struct {
	Code string `json:"code" binding:"required,len=6"`
}
