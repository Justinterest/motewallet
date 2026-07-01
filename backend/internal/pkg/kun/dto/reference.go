package dto

// CountriesReq is the body for POST /rest/v2.0/customer/fiat/withdrawal/countries.
// See: https://opendocs.kun.global/docs/api/get-countries-and-regions
// Currency is required when scene is WITHDRAWAL or BIND_ACCOUNT (USD/HKD/EUR).
type CountriesReq struct {
	RequestNo string `json:"requestNo"`
	Scene     string `json:"scene"`
	Currency  string `json:"currency,omitempty"`
	Language  string `json:"language,omitempty"`
}

// CountryItem is one entry in the countries API data array.
type CountryItem struct {
	CountryName string            `json:"countryName"`
	CountryCode string            `json:"countryCode"`
	Data        map[string]string `json:"data,omitempty"`
}

// CountryAuthTypesReq is the body for POST /rest/v2.0/customer/country/auth/types.
// See: https://opendocs.kun.global/docs/api/get-country-authentication-types
type CountryAuthTypesReq struct {
	CountryCode string `json:"countryCode"`
	RequestNo   string `json:"requestNo"`
}

// CountryAuthTypeItem is one entry in the country auth types API data array.
type CountryAuthTypeItem struct {
	DocName string            `json:"docName"`
	DocCode string            `json:"docCode"`
	Data    map[string]string `json:"data,omitempty"`
}
