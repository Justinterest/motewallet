package response

import "time"

type AdminAuthChallengeResp struct {
	Status     string         `json:"status"`
	TempToken  string         `json:"temp_token,omitempty"`
	TotpSecret string         `json:"totp_secret,omitempty"`
	TotpURI    string         `json:"totp_uri,omitempty"`
	Admin      *AdminInfoResp `json:"admin,omitempty"`
}

type AdminUserItemResp struct {
	ID          uint64     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	TotpEnabled bool       `json:"totp_enabled"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateAdminUserResp struct {
	User            AdminUserItemResp `json:"user"`
	InitialPassword string            `json:"initial_password,omitempty"`
}

type ResetAdminPasswordResp struct {
	NewPassword string `json:"new_password"`
}
