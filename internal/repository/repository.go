package repository

import (
	"gorm.io/gorm"
)

var (
	CountryRepo      *CountryRepository
	UserRepo         *UserRepository
	SampleSourceRepo *SampleSourceRepository
)

func InitRepositories(db *gorm.DB) {
	CountryRepo = NewCountryRepo(db)
	UserRepo = NewUserRepo(db)
	SampleSourceRepo = NewSampleSourceRepo(db)
}
