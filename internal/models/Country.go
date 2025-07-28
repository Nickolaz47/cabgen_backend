package models

type Country struct {
	Code  string `gorm:"primaryKey" json:"code"`
	Pt    string `gorm:"not null" json:"pt"`
	En    string `gorm:"not null" json:"en"`
	Es    string `gorm:"not null" json:"es"`
	Users []User `gorm:"foreignKey:CountryCode;references:Code"`
}
