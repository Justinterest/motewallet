package request

type AdminLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminChangePasswordReq struct {
	TempToken   string `json:"temp_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=64"`
}
