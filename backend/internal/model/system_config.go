package model

type SystemConfig struct {
	BaseModel
	ConfigKey   string  `gorm:"column:config_key;type:varchar(64);not null;uniqueIndex:uk_system_configs_key" json:"config_key"`
	ConfigValue string  `gorm:"column:config_value;type:text;not null" json:"config_value"`
	Description *string `gorm:"column:description;type:varchar(255)" json:"description,omitempty"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}
