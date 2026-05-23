package model

type FeeTemplate struct {
	SoftDeleteModel
	Name        string  `gorm:"column:name;type:varchar(64);not null" json:"name"`
	Description *string `gorm:"column:description;type:varchar(255)" json:"description,omitempty"`
	IsDefault   bool    `gorm:"column:is_default;type:tinyint(1);not null;default:0;index:idx_fee_templates_is_default" json:"is_default"`
}

func (FeeTemplate) TableName() string {
	return "fee_templates"
}
