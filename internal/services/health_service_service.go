package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HealthService interface {
	FindAll(ctx context.Context) ([]models.HealthServiceAdminTableResponse, error)
	FindAllActive(ctx context.Context) ([]models.HealthServiceFormResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*models.HealthServiceAdminTableResponse, error)
	FindByName(ctx context.Context, name string) ([]models.HealthServiceAdminTableResponse, error)
	Create(ctx context.Context, input models.HealthServiceCreateInput) (*models.HealthServiceAdminTableResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.HealthServiceUpdateInput) (*models.HealthServiceAdminTableResponse, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type healthServiceService struct {
	Repo        repositories.HealthServiceRepository
	CountryRepo repositories.CountryRepository
	Logger      *zap.Logger
}

func NewHealthServiceService(
	repo repositories.HealthServiceRepository,
	countryRepo repositories.CountryRepository,
	logger *zap.Logger,
) HealthService {
	return &healthServiceService{
		Repo: repo, CountryRepo: countryRepo, Logger: logger,
	}
}

func (s *healthServiceService) FindAll(ctx context.Context) (
	[]models.HealthServiceAdminTableResponse, error) {
	healthServices, err := s.Repo.GetHealthServices(ctx)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "FindAll", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponses := make([]models.HealthServiceAdminTableResponse,
		len(healthServices))
	for i, hs := range healthServices {
		tableResponses[i] = hs.ToAdminTableResponse()
	}

	return tableResponses, nil
}

func (s *healthServiceService) FindAllActive(ctx context.Context) (
	[]models.HealthServiceFormResponse, error) {
	healthServices, err := s.Repo.GetActiveHealthServices(ctx)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "FindAllActive", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	formServices := make([]models.HealthServiceFormResponse,
		len(healthServices))
	for i, hs := range healthServices {
		formServices[i] = hs.ToFormResponse()
	}

	return formServices, nil
}

func (s *healthServiceService) FindByID(ctx context.Context, ID uuid.UUID) (
	*models.HealthServiceAdminTableResponse, error) {
	healthService, err := s.Repo.GetHealthServiceByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "FindByID", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "FindByID", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponse := healthService.ToAdminTableResponse()
	return &tableResponse, nil
}

func (s *healthServiceService) FindByName(ctx context.Context, name string) (
	[]models.HealthServiceAdminTableResponse, error) {
	healthServices, err := s.Repo.GetHealthServicesByName(ctx, name)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "FindByName", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponses := make([]models.HealthServiceAdminTableResponse,
		len(healthServices))
	for i, hs := range healthServices {
		tableResponses[i] = hs.ToAdminTableResponse()
	}
	return tableResponses, nil
}

func (s *healthServiceService) Create(
	ctx context.Context, input models.HealthServiceCreateInput) (
	*models.HealthServiceAdminTableResponse, error) {
	country, err := s.CountryRepo.GetCountryByCode(ctx, input.CountryCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"HealthServiceService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrInvalidCountryCode
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"HealthServiceService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	healthService := models.HealthService{
		Name:         input.Name,
		Type:         models.HealthServiceType(input.Type),
		CountryID:    country.ID,
		Country:      *country,
		City:         input.City,
		Contactant:   input.Contactant,
		ContactEmail: input.ContactEmail,
		ContactPhone: input.ContactPhone,
		IsActive:     input.IsActive,
	}

	existingHealthService, err := s.Repo.GetHealthServiceDuplicate(
		ctx, healthService.Name, uuid.UUID{},
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "Create", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if existingHealthService != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "Create",
			logging.DatabaseConflictError, err,
		)...)
		return nil, ErrConflict
	}

	if err := s.Repo.CreateHealthService(ctx, &healthService); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "Create", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponse := healthService.ToAdminTableResponse()
	return &tableResponse, nil
}

func (s *healthServiceService) Update(
	ctx context.Context, ID uuid.UUID, input models.HealthServiceUpdateInput) (
	*models.HealthServiceAdminTableResponse, error) {
	existingHealthService, err := s.Repo.GetHealthServiceByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "Update",
			logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "Update",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if input.Name != nil {
		duplicate, err := s.Repo.GetHealthServiceDuplicate(ctx, *input.Name, ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"HealthServiceService", "Update", logging.DatabaseError, err,
			)...)
			return nil, ErrInternal
		}

		if duplicate != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"HealthServiceService", "Update",
				logging.DatabaseConflictError, err,
			)...)
			return nil, ErrConflict
		}
	}

	if input.CountryCode != nil {
		country, err := s.CountryRepo.GetCountryByCode(ctx, *input.CountryCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"HealthServiceService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrInvalidCountryCode
			}

			s.Logger.Error("Service Error", logging.ServiceLogging(
				"HealthServiceService", "Update",
				logging.ExternalRepositoryError, err,
			)...)
			return nil, ErrInternal
		}

		existingHealthService.CountryID = country.ID
		existingHealthService.Country = *country
	}

	validations.ApplyHealthServiceUpdate(existingHealthService, &input)

	if err := s.Repo.UpdateHealthService(ctx, existingHealthService); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"HealthServiceService", "Update",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	response := existingHealthService.ToAdminTableResponse()
	return &response, nil
}

func (s *healthServiceService) Delete(ctx context.Context, ID uuid.UUID) error {
	healthService, err := s.Repo.GetHealthServiceByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "Delete",
			logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if err := s.Repo.DeleteHealthService(ctx, healthService); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"HealthServiceService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	return nil
}
