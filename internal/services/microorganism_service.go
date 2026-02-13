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

type MicroorganismService interface {
	FindAll(ctx context.Context, language string) (
		[]models.MicroorganismAdminTableResponse, error)
	FindAllActive(ctx context.Context, language string) (
		[]models.MicroorganismFormResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (
		*models.MicroorganismAdminDetailResponse, error)
	FindBySpecies(ctx context.Context, species, language string) (
		[]models.MicroorganismAdminTableResponse, error)
	Create(ctx context.Context, input models.MicroorganismCreateInput) (
		*models.MicroorganismAdminDetailResponse, error)
	Update(ctx context.Context, ID uuid.UUID,
		input models.MicroorganismUpdateInput) (
		*models.MicroorganismAdminDetailResponse, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type microorganismService struct {
	Repo   repositories.MicroorganismRepository
	Logger *zap.Logger
}

func NewMicroorganismService(
	repo repositories.MicroorganismRepository,
	logger *zap.Logger,
) MicroorganismService {
	return &microorganismService{
		Repo:   repo,
		Logger: logger,
	}
}

func (s *microorganismService) FindAll(
	ctx context.Context, language string) (
	[]models.MicroorganismAdminTableResponse, error) {
	micros, err := s.Repo.GetMicroorganisms(ctx)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "FindAll",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	microResponses := make(
		[]models.MicroorganismAdminTableResponse, len(micros))
	for i, micro := range micros {
		microResponses[i] = micro.ToAdminTableResponse(language)
	}

	return microResponses, nil
}

func (s *microorganismService) FindAllActive(
	ctx context.Context, language string) (
	[]models.MicroorganismFormResponse, error) {
	micros, err := s.Repo.GetActiveMicroorganisms(ctx)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "FindAllActive",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	microResponses := make([]models.MicroorganismFormResponse, len(micros))
	for i, micro := range micros {
		microResponses[i] = micro.ToFormResponse(language)
	}

	return microResponses, nil
}

func (s *microorganismService) FindByID(
	ctx context.Context, ID uuid.UUID) (
	*models.MicroorganismAdminDetailResponse, error) {
	micro, err := s.Repo.GetMicroorganismByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "FindByID",
			logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "FindByID",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	detailResponse := micro.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *microorganismService) FindBySpecies(
	ctx context.Context, species, language string) (
	[]models.MicroorganismAdminTableResponse, error) {
	micros, err := s.Repo.GetMicroorganismsBySpecies(ctx, species, language)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "FindBySpecies",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	microResponses := make([]models.MicroorganismAdminTableResponse,
		len(micros))
	for i, micro := range micros {
		microResponses[i] = micro.ToAdminTableResponse(language)
	}

	return microResponses, nil
}

func (s *microorganismService) Create(
	ctx context.Context, input models.MicroorganismCreateInput) (
	*models.MicroorganismAdminDetailResponse, error) {
	micro := models.Microorganism{
		Taxon:    input.Taxon,
		Species:  input.Species,
		Variety:  input.Variety,
		IsActive: input.IsActive,
	}

	existingMicro, err := s.Repo.GetMicroorganismDuplicate(ctx, input.Species,
		input.Variety, uuid.Nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Create",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if existingMicro != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Create",
			logging.DatabaseConflictError, err,
		)...)
		return nil, ErrConflict
	}

	if err := s.Repo.CreateMicroorganism(ctx, &micro); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Create",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	detailedResponse := micro.ToAdminDetailResponse()
	return &detailedResponse, nil
}

func (s *microorganismService) Update(
	ctx context.Context, ID uuid.UUID, input models.MicroorganismUpdateInput) (
	*models.MicroorganismAdminDetailResponse, error) {
	existingMicro, err := s.Repo.GetMicroorganismByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Update",
			logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Update",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	validations.ApplyMicroorganismUpdate(existingMicro, &input)

	if input.Species != nil {
		duplicate, err := s.Repo.GetMicroorganismDuplicate(ctx,
			*input.Species, input.Variety, ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"MicroorganismService", "Update",
				logging.DatabaseError, err,
			)...)
			return nil, ErrInternal
		}

		if duplicate != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"MicroorganismService", "Update",
				logging.DatabaseConflictError, err,
			)...)
			return nil, ErrConflict
		}
	}

	if err := s.Repo.UpdateMicroorganism(ctx, existingMicro); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Update",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	detailResponse := existingMicro.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *microorganismService) Delete(ctx context.Context,
	ID uuid.UUID) error {
	micro, err := s.Repo.GetMicroorganismByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Delete",
			logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Delete",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if err := s.Repo.DeleteMicroorganism(ctx, micro); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MicroorganismService", "Delete",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	return nil
}
