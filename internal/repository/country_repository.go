package repository

import (
	"sync"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/gorm"
)

var (
	countryRepo     *CountryRepository
	countryRepoOnce sync.Once
)

type CountryRepository struct {
	DB *gorm.DB
}

func NewCountryRepo(db *gorm.DB) *CountryRepository {
	return &CountryRepository{DB: db}
}

func GetCountryRepo() *CountryRepository {
	countryRepoOnce.Do(func() {
		countryRepo = NewCountryRepo(db.DB)
	})
	return countryRepo
}

// Test only
func SetCountryRepo(r *CountryRepository) {
	countryRepo = r
}

func (r *CountryRepository) GetCountries() ([]models.Country, error) {
	var countries []models.Country
	if err := r.DB.Find(&countries).Error; err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *CountryRepository) GetCountry(contryCode string) (*models.Country, error) {
	var country models.Country
	if err := r.DB.Where("code = ?", contryCode).First(&country).Error; err != nil {
		return nil, err
	}

	return &country, nil
}
