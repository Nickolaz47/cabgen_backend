package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LaboratoryService interface {
	FindAll(ctx context.Context) ([]models.Laboratory, error)
	FindAllActive(ctx context.Context) ([]models.LaboratoryFormResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error)
	FindByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error)
	Create(ctx context.Context, lab *models.Laboratory) error
	Update(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type laboratoryService struct {
	Repo repository.LaboratoryRepository
}

func NewLaboratoryService(repo repository.LaboratoryRepository) LaboratoryService {
	return &laboratoryService{Repo: repo}
}

func (s *laboratoryService) FindAll(ctx context.Context) ([]models.Laboratory, error) {
	labs, err := s.Repo.GetLaboratories(ctx)

	if err != nil {
		return nil, ErrInternal
	}

	return labs, nil
}

func (s *laboratoryService) FindAllActive(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
	labs, err := s.Repo.GetActiveLaboratories(ctx)

	if err != nil {
		return nil, ErrInternal
	}

	formLabs := make([]models.LaboratoryFormResponse, len(labs))
	for i, lab := range labs {
		formLabs[i] = lab.ToFormResponse()
	}

	return formLabs, nil
}

func (s *laboratoryService) FindByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
	lab, err := s.Repo.GetLaboratoryByID(ctx, ID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	return lab, nil
}

func (s *laboratoryService) FindByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error) {
	labs, err := s.Repo.GetLaboratoriesByNameOrAbbreviation(ctx, input)

	if err != nil {
		return nil, ErrInternal
	}

	return labs, nil
}

func (s *laboratoryService) Create(ctx context.Context, lab *models.Laboratory) error {
	existingLab, err := s.Repo.GetLaboratoryDuplicate(ctx, lab.Name, uuid.UUID{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrInternal
	}

	if existingLab != nil {
		return ErrConflict
	}

	if err := s.Repo.CreateLaboratory(ctx, lab); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *laboratoryService) Update(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error) {
	existingLab, err := s.Repo.GetLaboratoryByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	
	if err != nil {
		return nil, ErrInternal
	}

	validations.ApplyLaboratoryUpdate(existingLab, &input)

	duplicate, err := s.Repo.GetLaboratoryDuplicate(ctx, existingLab.Name, existingLab.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}

	if duplicate != nil {
		return nil, ErrConflict
	}

	if err := s.Repo.UpdateLaboratory(ctx, existingLab); err != nil {
		return nil, ErrInternal
	}

	return existingLab, nil
}

func (s *laboratoryService) Delete(ctx context.Context, ID uuid.UUID) error {
	lab, err := s.Repo.GetLaboratoryByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return ErrInternal
	}

	if err := s.Repo.DeleteLaboratory(ctx, lab); err != nil {
		return ErrInternal
	}

	return nil
}
