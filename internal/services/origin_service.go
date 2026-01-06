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

type OriginService interface {
	FindAll(ctx context.Context) ([]models.Origin, error)
	FindAllActive(ctx context.Context, lang string) ([]models.OriginFormResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*models.Origin, error)
	FindByName(ctx context.Context, name, lang string) ([]models.Origin, error)
	Create(ctx context.Context, origin *models.Origin) error
	Update(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type originService struct {
	Repo repository.OriginRepository
}

func NewOriginService(repo repository.OriginRepository) OriginService {
	return &originService{Repo: repo}
}

func (s *originService) FindAll(ctx context.Context) ([]models.Origin, error) {
	origins, err := s.Repo.GetOrigins(ctx)

	if err != nil {
		return nil, ErrInternal
	}

	return origins, nil
}

func (s *originService) FindAllActive(ctx context.Context, lang string) ([]models.OriginFormResponse, error) {
	origins, err := s.Repo.GetActiveOrigins(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	formOrigins := make([]models.OriginFormResponse, len(origins))
	for i, origin := range origins {
		formOrigins[i] = origin.ToFormResponse(lang)
	}

	return formOrigins, nil
}

func (s *originService) FindByID(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
	origin, err := s.Repo.GetOriginByID(ctx, ID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	return origin, nil
}

func (s *originService) FindByName(ctx context.Context, name, lang string) ([]models.Origin, error) {
	origins, err := s.Repo.GetOriginsByName(ctx, name, lang)
	if err != nil {
		return nil, ErrInternal
	}

	return origins, nil
}

func (s *originService) Create(ctx context.Context, origin *models.Origin) error {
	existingOrigin, err := s.Repo.GetOriginDuplicate(ctx, origin.Names, uuid.UUID{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrInternal
	}

	if existingOrigin != nil {
		return ErrConflict
	}

	if err := s.Repo.CreateOrigin(ctx, origin); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *originService) Update(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error) {
	existingOrigin, err := s.Repo.GetOriginByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	validations.ApplyOriginUpdate(existingOrigin, &input)

	duplicate, err := s.Repo.GetOriginDuplicate(ctx, existingOrigin.Names, existingOrigin.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}

	if duplicate != nil {
		return nil, ErrConflict
	}

	if err := s.Repo.UpdateOrigin(ctx, existingOrigin); err != nil {
		return nil, ErrInternal
	}

	return existingOrigin, nil
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
