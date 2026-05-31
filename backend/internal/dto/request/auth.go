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
