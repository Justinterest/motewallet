package request

type SubmitKycReq struct {
	CompanyName        string `json:"company_name" binding:"required"`
	Country            string `json:"country" binding:"required"`
	RegistrationNumber string `json:"registration_number" binding:"required"`
	ContactName        string `json:"contact_name" binding:"required"`
	ContactPhone       string `json:"contact_phone" binding:"required"`
}
