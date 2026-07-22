package request

type CreateAdminUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email,max=128"`
	Role     string `json:"role" binding:"required,oneof=SUPER_ADMIN OPERATOR FINANCE"`
	Password string `json:"password" binding:"omitempty,min=8,max=64"`
}
