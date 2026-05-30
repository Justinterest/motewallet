package dto

// FileRef is a KUN-uploaded file reference (path from file upload API).
type FileRef struct {
	Path string `json:"path"`
}

// EnterpriseInfo is the enterprise section of sub-merchant onboarding authentication.
// See: https://opendocs.kun.global/docs/api/sub-merchant-onboarding-authentication
type EnterpriseInfo struct {
	IncorporationCertificate       []FileRef `json:"incorporationCertificate,omitempty"`
	IncorporationCertificateNo     string    `json:"incorporationCertificateNo"`
	EstablishTime                    string    `json:"establishTime"`
	EnterpriseEN                     string    `json:"enterpriseEN"`
	EnterpriseNameCHS                string    `json:"enterpriseNameCHS"`
	RegisterRegion                   string    `json:"registerRegion"`
	RegisterAddress                  string    `json:"registerAddress"`
	BusinessRegistration             []FileRef `json:"businessRegistration,omitempty"`
	BusinessRegistrationNo           string    `json:"businessRegistrationNo,omitempty"`
	Phone                            string    `json:"phone,omitempty"`
	IsChangeEnterpriseNameInFiveYears string `json:"isChangeEnterpriseNameInFiveYears,omitempty"`
	UsedEnterpriseName               string    `json:"usedEnterpriseName,omitempty"`
	EnterpriseType                   string    `json:"enterpriseType"`
	EnterpriseDomain                 string    `json:"enterpriseDomain,omitempty"`
	BusinessRegion                   []string  `json:"businessRegion,omitempty"`
	MainBusinessAddress              string    `json:"mainBusinessAddress,omitempty"`
	Industry                         string    `json:"industry"`
	SubIndustry                      string    `json:"subIndustry"`
	InitialFundingSource             string    `json:"initialFundingSource"`
	WealthSource                     string    `json:"wealthSource,omitempty"`
	ContinuousFundingSource          string    `json:"continuousFundingSource,omitempty"`
	SalesVolumeLastyear              string    `json:"salesVolumeLastyear,omitempty"`
	EmployeeNum                      string    `json:"employeeNum,omitempty"`
	OpenAccountPurpose               string    `json:"openAccountPurpose"`
	Incumbency                       []FileRef `json:"incumbency,omitempty"`
	AssociationRules                 []FileRef `json:"associationRules,omitempty"`
	AuthenticMaterials               []FileRef `json:"authenticMaterials,omitempty"`
	ManagerCountry                   string    `json:"managerCountry"`
	ManagerAuthType                  string    `json:"managerAuthType"`
	ManagerVerificationType          string    `json:"managerVerificationType"`
	ManagerIdHolding                 []FileRef `json:"managerIdHolding,omitempty"`
	ManagerCertificateTerm           []string  `json:"managerCertificateTerm,omitempty"`
	ManagerSurnameCHS                string    `json:"managerSurnameCHS,omitempty"`
	ManagerNameCHS                   string    `json:"managerNameCHS,omitempty"`
	ManagerSurname                   string    `json:"managerSurname,omitempty"`
	ManagerNameEN                    string    `json:"managerNameEN"`
	ManagerBirthday                  string    `json:"managerBirthday"`
	ManagerGender                    string    `json:"managerGender"`
	ManagerIdCard                    string    `json:"managerIdCard"`
	ManagerResidenceCountry          string    `json:"managerResidenceCountry"`
	ManagerResidenceAddress          string    `json:"managerResidenceAddress"`
	ManagerContactsEmail             string    `json:"managerContactsEmail,omitempty"`
	AuthorizationLetter              []FileRef `json:"authorizationLetter,omitempty"`
	EquityStructure                  []FileRef `json:"equityStructure,omitempty"`
	MiddleTierShareholders           string    `json:"middleTierShareholders,omitempty"`
	Nnc1                             []FileRef `json:"nnc1,omitempty"`
}

// PersonInfo is shared by shareholders and directors.
type PersonInfo struct {
	Gender              string      `json:"gender"`
	IdCard              string      `json:"idCard"`
	NameEN              string      `json:"nameEN"`
	Country             string      `json:"country"`
	NameCHS             string      `json:"nameCHS,omitempty"`
	Surname             string      `json:"surname"`
	AuthType            string      `json:"authType"`
	Birthday            string      `json:"birthday"`
	SurnameCHS          string      `json:"surnameCHS,omitempty"`
	CertificateTerm     interface{} `json:"certificateTerm,omitempty"` // []string or object per KUN
	ResidenceAddress    string      `json:"residenceAddress"`
	ResidenceCountry    string      `json:"residenceCountry"`
	ShareholdingRatio   string      `json:"shareholdingRatio,omitempty"`
	IdHolding           []FileRef   `json:"idHolding,omitempty"`
	VerificationType    string      `json:"verificationType,omitempty"`
}

// SubMerchantRegisterReq is the request body for POST /rest/v2.0/customer/subMerchant/register.
// Customer-No header must be the sub-merchant number (not the platform customer no).
type SubMerchantRegisterReq struct {
	RequestNo        string         `json:"requestNo"`
	EnterpriseInfo   EnterpriseInfo `json:"enterpriseInfo"`
	ShareholdersInfo []PersonInfo   `json:"shareholdersInfo"`
	DirectorInfo     []PersonInfo   `json:"directorInfo"`
}

type SubMerchantRegisterResp struct {
	AuthID string   `json:"authId"`
	Errors []string `json:"errors,omitempty"`
}

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

type MerchantRegisterQueryReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
}

type MerchantRegisterQueryResp struct {
	AuthStatus string `json:"authStatus"`
	FailReason string `json:"failReason,omitempty"`
}
