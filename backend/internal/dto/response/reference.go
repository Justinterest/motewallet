package response

// CountryOption is a country/region choice for KYC forms.
type CountryOption struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
}

// CountryListResp lists countries for a given scene.
type CountryListResp struct {
	Items []CountryOption `json:"items"`
}

// AuthTypeOption is a document type choice for a country.
type AuthTypeOption struct {
	DocCode string `json:"doc_code"`
	DocName string `json:"doc_name"`
}

// CountryAuthTypesResp lists authentication document types for a country.
type CountryAuthTypesResp struct {
	Items []AuthTypeOption `json:"items"`
}
