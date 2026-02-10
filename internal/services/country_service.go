package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"go.uber.org/zap"
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
	Repo   repositories.CountryRepository
	Logger *zap.Logger
}

func NewCountryService(repo repositories.CountryRepository,
	logger *zap.Logger) CountryService {
	return &countryService{Repo: repo, Logger: logger}
}

func (s *countryService) FindAll(ctx context.Context, language string) ([]models.CountryFormResponse, error) {
	countries, err := s.Repo.GetCountries(ctx)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "FindAll", logging.DatabaseError, err,
		)...)
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
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "FindByCode", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "FindByCode", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	detailResponse := country.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *countryService) FindByName(ctx context.Context, name, language string) ([]models.CountryFormResponse, error) {
	countries, err := s.Repo.GetCountriesByName(ctx, name, language)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "FindByName", logging.DatabaseError, err,
		)...)
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
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Create", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if existingCountry != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Create", logging.DatabaseConflictError, err,
		)...)
		return nil, ErrConflict
	}

	if err := s.Repo.CreateCountry(ctx, &country); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Create", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	response := country.ToAdminDetailResponse()
	return &response, nil
}

func (s *countryService) Update(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error) {
	existingCountry, err := s.Repo.GetCountryByCode(ctx, code)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Update", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Update", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	validations.ApplyCountryUpdate(existingCountry, &input)

	if input.Code != nil && *input.Code != existingCountry.Code {
		_, err := s.Repo.GetCountryByCode(ctx, *input.Code)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"CountryService", "Update", logging.DatabaseError, err,
			)...)
			return nil, ErrInternal
		}

		if err == nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"CountryService", "Update", logging.DatabaseConflictError, err,
			)...)
			return nil, ErrConflict
		}
	}

	if input.Names != nil {
		duplicate, err := s.Repo.GetCountryDuplicate(ctx, input.Names, code)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"CountryService", "Update", logging.DatabaseError, err,
			)...)
			return nil, ErrInternal
		}

		if duplicate != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"CountryService", "Update", logging.DatabaseConflictError, err,
			)...)
			return nil, ErrConflict
		}
	}

	if err := s.Repo.UpdateCountry(ctx, existingCountry); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Update", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	response := existingCountry.ToAdminDetailResponse()
	return &response, nil
}

func (s *countryService) Delete(ctx context.Context, code string) error {
	country, err := s.Repo.GetCountryByCode(ctx, code)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Delete", logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if err := s.Repo.DeleteCountry(ctx, country); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"CountryService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	return nil
}
