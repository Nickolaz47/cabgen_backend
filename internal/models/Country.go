package models

type Country struct {
	Code  string `gorm:"primaryKey"`
	Pt    string `gorm:"not null"`
	En    string `gorm:"not null"`
	Es    string `gorm:"not null"`
	Users []User `gorm:"foreignKey:CountryCode;references:Code"`
}
