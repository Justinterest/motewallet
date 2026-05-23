package dto

type DepositAddressReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	Currency      string `json:"currency"`
	Chain         string `json:"chain"`
}

type DepositAddressResp struct {
	Address  string `json:"address"`
	Chain    string `json:"chain"`
	Currency string `json:"currency"`
}

type BalanceQueryReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	Currency      string `json:"currency"`
	CurrencyType  string `json:"currencyType,omitempty"`
}

type BalanceQueryResp struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Total     string `json:"total"`
}

type CryptoAddressAddReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	Currency      string `json:"currency"`
	Chain         string `json:"chain"`
	Address       string `json:"address"`
	Alias         string `json:"alias,omitempty"`
	RequestNo     string `json:"requestNo"`
}

type CryptoAddressAddResp struct {
	AccountId string `json:"accountId"`
}

type FiatAddressAddReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	Currency      string `json:"currency"`
	BankName      string `json:"bankName"`
	BankCountry   string `json:"bankCountry"`
	SwiftCode     string `json:"swiftCode"`
	AccountName   string `json:"accountName"`
	AccountNo     string `json:"accountNo"`
	TransferType  string `json:"transferType"`
	RequestNo     string `json:"requestNo"`
}

type FiatAddressAddResp struct {
	AccountId string `json:"accountId"`
}

type RegionBalanceQueryReq struct {
	SubCustomerNo string `json:"subCustomerNo"`
	RegionCode    string `json:"regionCode"`
}

type RegionBalanceItem struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Total     string `json:"total"`
}

type RegionBalanceQueryResp struct {
	Balances []RegionBalanceItem `json:"balances"`
}
