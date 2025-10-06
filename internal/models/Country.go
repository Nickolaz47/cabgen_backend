package models

type Country struct {
	Code  string  `gorm:"primaryKey" json:"code"`
	Names JSONMap `gorm:"type:jsonb;not null" json:"names"`
	Users []User  `gorm:"foreignKey:CountryCode;references:Code" json:"-"`
}
