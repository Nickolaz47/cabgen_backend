package models

type Country struct {
	Code  string            `gorm:"primaryKey" json:"code"`
	Names map[string]string `gorm:"json;not null" json:"names"`
	Users []User            `gorm:"foreignKey:CountryCode;references:Code" json:"-"`
}
