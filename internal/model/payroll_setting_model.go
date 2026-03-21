package model

import "time"

type PayrollSetting struct {
	ID         string    `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	ConfigName string    `gorm:"column:config_name;type:varchar(255);not null"`
	ConfigKey  string    `gorm:"column:config_key;type:varchar(255);unique;not null"`
	Value      float64   `gorm:"column:value;type:decimal(12,2);not null;default:0"`
	IsACtive   bool      `gorm:"column:is_active;type:boolean;not null;default:true"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (PayrollSetting) TableName() string {
	return "payroll_settings"
}
