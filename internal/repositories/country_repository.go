package repositories

import (
	"context"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/gorm"
)

type CountryRepository interface {
	GetCountries(ctx context.Context) ([]models.Country, error)
	GetCountryByID(ctx context.Context, ID uint) (*models.Country, error)
	GetCountryByCode(ctx context.Context, code string) (*models.Country, error)
	GetCountriesByName(ctx context.Context, name, lang string) ([]models.Country, error)
	GetCountryDuplicate(ctx context.Context, names models.JSONMap, code string) (*models.Country, error)
	CreateCountry(ctx context.Context, country *models.Country) error
	UpdateCountry(ctx context.Context, country *models.Country) error
	DeleteCountry(ctx context.Context, country *models.Country) error
}

type countryRepo struct {
	DB *gorm.DB
}

func NewCountryRepo(db *gorm.DB) CountryRepository {
	return &countryRepo{DB: db}
}

func (r *countryRepo) GetCountries(ctx context.Context) ([]models.Country, error) {
	var countries []models.Country
	if err := r.DB.WithContext(ctx).Find(&countries).Error; err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *countryRepo) GetCountryByID(ctx context.Context, ID uint) (*models.Country, error) {
	var country models.Country
	if err := r.DB.WithContext(ctx).Where("id = ?", ID).First(&country).Error; err != nil {
		return nil, err
	}

	return &country, nil
}

func (r *countryRepo) GetCountryByCode(ctx context.Context, code string) (*models.Country, error) {
	var country models.Country
	if err := r.DB.WithContext(ctx).Where("code = ?", code).First(&country).Error; err != nil {
		return nil, err
	}

	return &country, nil
}

func (r *countryRepo) GetCountriesByName(ctx context.Context, name, lang string) ([]models.Country, error) {
	var countries []models.Country
	query := "LOWER(names->>'" + lang + "') LIKE LOWER(?)"
	if err := r.DB.WithContext(ctx).Where(query, "%"+name+"%").Find(&countries).Error; err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *countryRepo) GetCountryDuplicate(ctx context.Context, names models.JSONMap, code string) (*models.Country, error) {
	var country models.Country

	conditions := r.DB.WithContext(ctx)
	for lang, value := range names {
		conditions = conditions.Or(
			fmt.Sprintf(
				"LOWER(names->>'%s') = LOWER(?)",
				lang,
			),
			value,
		)
	}

	query := r.DB.WithContext(ctx).Where(conditions)

	if code != "" {
		query = query.Where("code != ?", code)
	}

	if err := query.First(&country).Error; err != nil {
		return nil, err
	}

	return &country, nil
}

func (r *countryRepo) CreateCountry(ctx context.Context, country *models.Country) error {
	return r.DB.WithContext(ctx).Create(country).Error
}

func (r *countryRepo) UpdateCountry(ctx context.Context, country *models.Country) error {
	return r.DB.WithContext(ctx).Save(country).Error
}

func (r *countryRepo) DeleteCountry(ctx context.Context, country *models.Country) error {
	return r.DB.WithContext(ctx).Delete(country).Error
}
