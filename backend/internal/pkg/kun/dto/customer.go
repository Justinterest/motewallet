package dto

type RegisterReq struct {
	Email     string `json:"email"`
	RequestNo string `json:"requestNo"`
}

type RegisterResp struct {
	SubCustomerNo string `json:"subCustomerNo"`
}

type AgreeListReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	SignStatus    string `json:"signStatus"`
	BizCode       string `json:"bizCode"`
}

type Agreement struct {
	ProtocolId string `json:"protocolId"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	SignStatus string `json:"signStatus"`
}

type AgreeListResp struct {
	List []Agreement `json:"list"`
}

type AgreeAuthReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	ProtocolIds   string `json:"protocolIds"`
}

type SubMerchantRegisterReq struct {
	SubCustomerNo      string `json:"subCustomerNo"`
	CompanyName        string `json:"companyName"`
	Country            string `json:"country"`
	RegistrationNumber string `json:"registrationNumber"`
	BusinessType       string `json:"businessType,omitempty"`
	ContactName        string `json:"contactName,omitempty"`
	ContactEmail       string `json:"contactEmail,omitempty"`
	ContactPhone       string `json:"contactPhone,omitempty"`
	RequestNo          string `json:"requestNo"`
}

type MerchantRegisterQueryReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
}

type MerchantRegisterQueryResp struct {
	AuthStatus string `json:"authStatus"`
	FailReason string `json:"failReason,omitempty"`
}
