package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OriginService interface {
	FindAll(ctx context.Context, language string) ([]models.OriginAdminTableResponse, error)
	FindAllActive(ctx context.Context, language string) ([]models.OriginFormResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*models.OriginAdminDetailResponse, error)
	FindByName(ctx context.Context, name, language string) ([]models.OriginAdminTableResponse, error)
	Create(ctx context.Context, input models.OriginCreateInput) (*models.OriginAdminDetailResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.OriginAdminDetailResponse, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type originService struct {
	Repo repositories.OriginRepository
}

func NewOriginService(repo repositories.OriginRepository) OriginService {
	return &originService{Repo: repo}
}

func (s *originService) FindAll(ctx context.Context, language string) ([]models.OriginAdminTableResponse, error) {
	origins, err := s.Repo.GetOrigins(ctx)

	if err != nil {
		return nil, ErrInternal
	}

	tableResponses := make([]models.OriginAdminTableResponse, len(origins))
	for i, origin := range origins {
		tableResponses[i] = origin.ToAdminTableResponse(language)
	}

	return tableResponses, nil
}

func (s *originService) FindAllActive(ctx context.Context, language string) ([]models.OriginFormResponse, error) {
	origins, err := s.Repo.GetActiveOrigins(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	formOrigins := make([]models.OriginFormResponse, len(origins))
	for i, origin := range origins {
		formOrigins[i] = origin.ToFormResponse(language)
	}

	return formOrigins, nil
}

func (s *originService) FindByID(ctx context.Context, ID uuid.UUID) (*models.OriginAdminDetailResponse, error) {
	origin, err := s.Repo.GetOriginByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	detailResponse := origin.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *originService) FindByName(ctx context.Context, name, language string) ([]models.OriginAdminTableResponse, error) {
	origins, err := s.Repo.GetOriginsByName(ctx, name, language)
	if err != nil {
		return nil, ErrInternal
	}

	tableResponses := make([]models.OriginAdminTableResponse, len(origins))
	for i, origin := range origins {
		tableResponses[i] = origin.ToAdminTableResponse(language)
	}

	return tableResponses, nil
}

func (s *originService) Create(ctx context.Context, input models.OriginCreateInput) (*models.OriginAdminDetailResponse, error) {
	origin := models.Origin{
		Names:    input.Names,
		IsActive: input.IsActive,
	}

	existingOrigin, err := s.Repo.GetOriginDuplicate(ctx, origin.Names, uuid.UUID{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}

	if existingOrigin != nil {
		return nil, ErrConflict
	}

	if err := s.Repo.CreateOrigin(ctx, &origin); err != nil {
		return nil, ErrInternal
	}

	detailResponse := origin.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *originService) Update(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.OriginAdminDetailResponse, error) {
	existingOrigin, err := s.Repo.GetOriginByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	validations.ApplyOriginUpdate(existingOrigin, &input)

	if input.Names != nil {
		duplicate, err := s.Repo.GetOriginDuplicate(ctx, input.Names, ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInternal
		}

		if duplicate != nil {
			return nil, ErrConflict
		}
	}

	if err := s.Repo.UpdateOrigin(ctx, existingOrigin); err != nil {
		return nil, ErrInternal
	}

	detailResponse := existingOrigin.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *originService) Delete(ctx context.Context, ID uuid.UUID) error {
	origin, err := s.Repo.GetOriginByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return ErrInternal
	}

	if err := s.Repo.DeleteOrigin(ctx, origin); err != nil {
		return ErrInternal
	}

	return nil
}
