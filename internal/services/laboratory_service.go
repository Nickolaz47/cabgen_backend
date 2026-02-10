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

type LaboratoryService interface {
	FindAll(ctx context.Context) ([]models.LaboratoryAdminTableResponse, error)
	FindAllActive(ctx context.Context) ([]models.LaboratoryFormResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*models.LaboratoryAdminTableResponse, error)
	FindByNameOrAbbreviation(ctx context.Context, input string) ([]models.LaboratoryAdminTableResponse, error)
	Create(ctx context.Context, input models.LaboratoryCreateInput) (*models.LaboratoryAdminTableResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.LaboratoryAdminTableResponse, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type laboratoryService struct {
	Repo   repositories.LaboratoryRepository
	Logger *zap.Logger
}

func NewLaboratoryService(repo repositories.LaboratoryRepository,
	logger *zap.Logger) LaboratoryService {
	return &laboratoryService{Repo: repo, Logger: logger}
}

func (s *laboratoryService) FindAll(ctx context.Context) ([]models.LaboratoryAdminTableResponse, error) {
	labs, err := s.Repo.GetLaboratories(ctx)

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "FindAll", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponses := make([]models.LaboratoryAdminTableResponse, len(labs))
	for i, lab := range labs {
		tableResponses[i] = lab.ToAdminTableResponse()
	}
	return tableResponses, nil
}

func (s *laboratoryService) FindAllActive(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
	labs, err := s.Repo.GetActiveLaboratories(ctx)

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "FindAllActive", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	formLabs := make([]models.LaboratoryFormResponse, len(labs))
	for i, lab := range labs {
		formLabs[i] = lab.ToFormResponse()
	}

	return formLabs, nil
}

func (s *laboratoryService) FindByID(ctx context.Context, ID uuid.UUID) (*models.LaboratoryAdminTableResponse, error) {
	lab, err := s.Repo.GetLaboratoryByID(ctx, ID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "FindByID", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "FindByID", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponse := lab.ToAdminTableResponse()
	return &tableResponse, nil
}

func (s *laboratoryService) FindByNameOrAbbreviation(
	ctx context.Context,
	input string) ([]models.LaboratoryAdminTableResponse, error) {
	labs, err := s.Repo.GetLaboratoriesByNameOrAbbreviation(ctx, input)

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "FindByNameOrAbbreviation",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponses := make([]models.LaboratoryAdminTableResponse, len(labs))
	for i, lab := range labs {
		tableResponses[i] = lab.ToAdminTableResponse()
	}
	return tableResponses, nil
}

func (s *laboratoryService) Create(
	ctx context.Context,
	input models.LaboratoryCreateInput) (*models.LaboratoryAdminTableResponse, error) {
	lab := models.Laboratory{
		Name:         input.Name,
		Abbreviation: input.Abbreviation,
		IsActive:     input.IsActive,
	}

	existingLab, err := s.Repo.GetLaboratoryDuplicate(ctx, lab.Name, uuid.UUID{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Create",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if existingLab != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Create",
			logging.DatabaseConflictError, err,
		)...)
		return nil, ErrConflict
	}

	if err := s.Repo.CreateLaboratory(ctx, &lab); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Create",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponse := lab.ToAdminTableResponse()
	return &tableResponse, nil
}

func (s *laboratoryService) Update(
	ctx context.Context,
	ID uuid.UUID,
	input models.LaboratoryUpdateInput) (*models.LaboratoryAdminTableResponse, error) {
	existingLab, err := s.Repo.GetLaboratoryByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Update",
			logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Update",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	validations.ApplyLaboratoryUpdate(existingLab, &input)

	if input.Name != nil {
		duplicate, err := s.Repo.GetLaboratoryDuplicate(ctx, *input.Name, ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"LaboratoryService", "Update",
				logging.DatabaseError, err,
			)...)
			return nil, ErrInternal
		}

		if duplicate != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"LaboratoryService", "Update",
				logging.DatabaseConflictError, err,
			)...)
			return nil, ErrConflict
		}
	}

	if err := s.Repo.UpdateLaboratory(ctx, existingLab); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Update",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	tableResponse := existingLab.ToAdminTableResponse()
	return &tableResponse, nil
}

func (s *laboratoryService) Delete(ctx context.Context, ID uuid.UUID) error {
	lab, err := s.Repo.GetLaboratoryByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Delete",
			logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Delete",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if err := s.Repo.DeleteLaboratory(ctx, lab); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"LaboratoryService", "Delete",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	return nil
}
