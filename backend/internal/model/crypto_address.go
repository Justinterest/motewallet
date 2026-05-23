package model

type CryptoAddress struct {
	SoftDeleteModel
	MerchantID   uint64  `gorm:"column:merchant_id;not null;index:idx_crypto_addresses_merchant_id;index:idx_crypto_addresses_merchant_currency" json:"merchant_id"`
	Currency     string  `gorm:"column:currency;type:varchar(10);not null;index:idx_crypto_addresses_merchant_currency" json:"currency"`
	Chain        string  `gorm:"column:chain;type:varchar(20);not null;index:idx_crypto_addresses_merchant_currency" json:"chain"`
	Address      string  `gorm:"column:address;type:varchar(255);not null" json:"address"`
	Alias        string  `gorm:"column:alias;type:varchar(64);not null" json:"alias"`
	KunAccountID *string `gorm:"column:kun_account_id;type:varchar(128)" json:"kun_account_id,omitempty"`
	Status       string  `gorm:"column:status;type:varchar(20);not null;default:INIT;index:idx_crypto_addresses_status" json:"status"`
}

func (CryptoAddress) TableName() string {
	return "crypto_addresses"
}
