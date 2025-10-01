package repository

import (
	"gorm.io/gorm"
)

var (
	CountryRepo *CountryRepository
	UserRepo    *UserRepository
	OriginRepo  *OriginRepository
)

func InitRepositories(db *gorm.DB) {
	CountryRepo = NewCountryRepo(db)
	UserRepo = NewUserRepo(db)
	OriginRepo = NewOriginRepo(db)
}
