package model

type BankAccount struct {
	SoftDeleteModel
	MerchantID       uint64  `gorm:"column:merchant_id;not null;index:idx_bank_accounts_merchant_id" json:"merchant_id"`
	KunAccountID     *string `gorm:"column:kun_account_id;type:varchar(128)" json:"kun_account_id,omitempty"`
	CurrencyList     string  `gorm:"column:currency_list;type:varchar(50);not null" json:"currency_list"`
	TransferType     string  `gorm:"column:transfer_type;type:varchar(10);not null" json:"transfer_type"`
	AccountNo        string  `gorm:"column:account_no;type:varchar(128);not null" json:"account_no"`
	AccountName      string  `gorm:"column:account_name;type:varchar(128);not null" json:"account_name"`
	BankName         *string `gorm:"column:bank_name;type:varchar(128)" json:"bank_name,omitempty"`
	BankCode         *string `gorm:"column:bank_code;type:varchar(20)" json:"bank_code,omitempty"`
	SwiftCode        *string `gorm:"column:swift_code;type:varchar(20)" json:"swift_code,omitempty"`
	PayeeCountryCode *string `gorm:"column:payee_country_code;type:varchar(10)" json:"payee_country_code,omitempty"`
	PayeeAddress     *string `gorm:"column:payee_address;type:varchar(255)" json:"payee_address,omitempty"`
	MiddleSwiftCode  *string `gorm:"column:middle_swift_code;type:varchar(20)" json:"middle_swift_code,omitempty"`
	Area             string  `gorm:"column:area;type:varchar(10);not null" json:"area"`
	Status           string  `gorm:"column:status;type:varchar(20);not null;default:ACTIVE" json:"status"`
}

func (BankAccount) TableName() string {
	return "bank_accounts"
}
