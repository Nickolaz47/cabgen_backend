package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type MockCountryRepository struct {
	GetCountriesFunc        func(ctx context.Context) ([]models.Country, error)
	GetCountryByIDFunc      func(ctx context.Context, ID uint) (*models.Country, error)
	GetCountryByCodeFunc    func(ctx context.Context, code string) (*models.Country, error)
	GetCountriesByNameFunc  func(ctx context.Context, name, lang string) ([]models.Country, error)
	GetCountryDuplicateFunc func(ctx context.Context, names models.JSONMap, code string) (*models.Country, error)
	CreateCountryFunc       func(ctx context.Context, country *models.Country) error
	UpdateCountryFunc       func(ctx context.Context, country *models.Country) error
	DeleteCountryFunc       func(ctx context.Context, country *models.Country) error
}

func (r *MockCountryRepository) GetCountries(ctx context.Context) ([]models.Country, error) {
	if r.GetCountriesFunc != nil {
		return r.GetCountriesFunc(ctx)
	}
	return nil, nil
}

func (r *MockCountryRepository) GetCountryByID(ctx context.Context, ID uint) (*models.Country, error) {
	if r.GetCountryByIDFunc != nil {
		return r.GetCountryByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (r *MockCountryRepository) GetCountryByCode(ctx context.Context, code string) (*models.Country, error) {
	if r.GetCountryByCodeFunc != nil {
		return r.GetCountryByCodeFunc(ctx, code)
	}
	return nil, nil
}

func (r *MockCountryRepository) GetCountriesByName(ctx context.Context, name, lang string) ([]models.Country, error) {
	if r.GetCountriesByNameFunc != nil {
		return r.GetCountriesByNameFunc(ctx, name, lang)
	}
	return nil, nil
}

func (r *MockCountryRepository) GetCountryDuplicate(ctx context.Context, names models.JSONMap, code string) (*models.Country, error) {
	if r.GetCountryDuplicateFunc != nil {
		return r.GetCountryDuplicateFunc(ctx, names, code)
	}
	return nil, nil
}

func (r *MockCountryRepository) CreateCountry(ctx context.Context, country *models.Country) error {
	if r.CreateCountryFunc != nil {
		return r.CreateCountryFunc(ctx, country)
	}
	return nil
}

func (r *MockCountryRepository) UpdateCountry(ctx context.Context, country *models.Country) error {
	if r.UpdateCountryFunc != nil {
		return r.UpdateCountryFunc(ctx, country)
	}
	return nil
}

func (r *MockCountryRepository) DeleteCountry(ctx context.Context, country *models.Country) error {
	if r.DeleteCountryFunc != nil {
		return r.DeleteCountryFunc(ctx, country)
	}
	return nil
}

type MockCountryService struct {
	FindAllFunc    func(ctx context.Context, lang string) ([]models.CountryFormResponse, error)
	FindByCodeFunc func(ctx context.Context, code string) (*models.CountryAdminDetailResponse, error)
	FindByNameFunc func(ctx context.Context, name, lang string) ([]models.CountryFormResponse, error)
	CreateFunc     func(ctx context.Context, input models.CountryCreateInput) (*models.CountryAdminDetailResponse, error)
	UpdateFunc     func(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error)
	DeleteFunc     func(ctx context.Context, code string) error
}

func (m *MockCountryService) FindAll(ctx context.Context, lang string) ([]models.CountryFormResponse, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, lang)
	}
	return nil, nil
}

func (m *MockCountryService) FindByCode(ctx context.Context, code string) (*models.CountryAdminDetailResponse, error) {
	if m.FindByCodeFunc != nil {
		return m.FindByCodeFunc(ctx, code)
	}
	return nil, nil
}

func (m *MockCountryService) FindByName(ctx context.Context, name, lang string) ([]models.CountryFormResponse, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(ctx, name, lang)
	}
	return nil, nil
}

func (m *MockCountryService) Create(ctx context.Context, input models.CountryCreateInput) (*models.CountryAdminDetailResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockCountryService) Update(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, code, input)
	}
	return nil, nil
}

func (m *MockCountryService) Delete(ctx context.Context, code string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, code)
	}
	return nil
}
