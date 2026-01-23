package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"gorm.io/gorm"
)

type CountryService interface {
	FindAll(ctx context.Context, language string) ([]models.CountryFormResponse, error)
	FindByCode(ctx context.Context, code string) (*models.CountryAdminDetailResponse, error)
	FindByName(ctx context.Context, name, language string) ([]models.CountryFormResponse, error)
	Create(ctx context.Context, input models.CountryCreateInput) (*models.CountryAdminDetailResponse, error)
	Update(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error)
	Delete(ctx context.Context, code string) error
}

type countryService struct {
	Repo repository.CountryRepository
}

func NewCountryService(repo repository.CountryRepository) CountryService {
	return &countryService{Repo: repo}
}

func (s *countryService) FindAll(ctx context.Context, language string) ([]models.CountryFormResponse, error) {
	countries, err := s.Repo.GetCountries(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	responses := make([]models.CountryFormResponse, len(countries))
	for i, country := range countries {
		responses[i] = country.ToFormResponse(language)
	}
	return responses, nil
}

func (s *countryService) FindByCode(ctx context.Context, code string) (*models.CountryAdminDetailResponse, error) {
	country, err := s.Repo.GetCountryByCode(ctx, code)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	detailResponse := country.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *countryService) FindByName(ctx context.Context, name, language string) ([]models.CountryFormResponse, error) {
	countries, err := s.Repo.GetCountriesByName(ctx, name, language)
	if err != nil {
		return nil, ErrInternal
	}

	responses := make([]models.CountryFormResponse, len(countries))
	for i, country := range countries {
		responses[i] = country.ToFormResponse(language)
	}
	return responses, nil
}

func (s *countryService) Create(ctx context.Context, input models.CountryCreateInput) (*models.CountryAdminDetailResponse, error) {
	country := models.Country{
		Code:  input.Code,
		Names: input.Names,
	}

	existingCountry, err := s.Repo.GetCountryDuplicate(ctx, input.Names, "")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}

	if existingCountry != nil {
		return nil, ErrConflict
	}

	if err := s.Repo.CreateCountry(ctx, &country); err != nil {
		return nil, ErrInternal
	}

	response := country.ToAdminDetailResponse()
	return &response, nil
}

func (s *countryService) Update(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error) {
	existingCountry, err := s.Repo.GetCountryByCode(ctx, code)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	validations.ApplyCountryUpdate(existingCountry, &input)

	if input.Code != nil && *input.Code != existingCountry.Code {
		_, err := s.Repo.GetCountryByCode(ctx, *input.Code)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInternal
		}

		if err == nil {
			return nil, ErrConflict
		}
	}

	if input.Names != nil {
		duplicate, err := s.Repo.GetCountryDuplicate(ctx, input.Names, code)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInternal
		}

		if duplicate != nil {
			return nil, ErrConflict
		}
	}

	if err := s.Repo.UpdateCountry(ctx, existingCountry); err != nil {
		return nil, ErrInternal
	}

	response := existingCountry.ToAdminDetailResponse()
	return &response, nil
}

func (s *countryService) Delete(ctx context.Context, code string) error {
	country, err := s.Repo.GetCountryByCode(ctx, code)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return ErrInternal
	}

	if err := s.Repo.DeleteCountry(ctx, country); err != nil {
		return ErrInternal
	}

	return nil
}
