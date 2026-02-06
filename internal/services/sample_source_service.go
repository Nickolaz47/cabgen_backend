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

type SampleSourceService interface {
	FindAll(ctx context.Context, language string) ([]models.SampleSourceAdminTableResponse, error)
	FindAllActive(ctx context.Context, language string) ([]models.SampleSourceFormResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*models.SampleSourceAdminDetailResponse, error)
	FindByNameOrGroup(ctx context.Context, input, language string) ([]models.SampleSourceAdminTableResponse, error)
	Create(ctx context.Context, input models.SampleSourceCreateInput) (*models.SampleSourceAdminDetailResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.SampleSourceUpdateInput) (*models.SampleSourceAdminDetailResponse, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type sampleSourceService struct {
	Repo repositories.SampleSourceRepository
}

func NewSampleSourceService(repo repositories.SampleSourceRepository) SampleSourceService {
	return &sampleSourceService{Repo: repo}
}

func (s *sampleSourceService) FindAll(ctx context.Context, language string) ([]models.SampleSourceAdminTableResponse, error) {
	sampleSources, err := s.Repo.GetSampleSources(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	tableResponses := make([]models.SampleSourceAdminTableResponse, len(sampleSources))
	for i, sample := range sampleSources {
		tableResponses[i] = sample.ToAdminTableResponse(language)
	}
	return tableResponses, nil
}

func (s *sampleSourceService) FindAllActive(ctx context.Context, language string) ([]models.SampleSourceFormResponse, error) {
	sampleSources, err := s.Repo.GetActiveSampleSources(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	formSampleSources := make([]models.SampleSourceFormResponse, len(sampleSources))
	for i, sample := range sampleSources {
		formSampleSources[i] = sample.ToFormResponse(language)
	}

	return formSampleSources, nil
}

func (s *sampleSourceService) FindByID(ctx context.Context, ID uuid.UUID) (*models.SampleSourceAdminDetailResponse, error) {
	sampleSource, err := s.Repo.GetSampleSourceByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	detailResponse := sampleSource.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *sampleSourceService) FindByNameOrGroup(ctx context.Context, input, language string) ([]models.SampleSourceAdminTableResponse, error) {
	sampleSources, err := s.Repo.GetSampleSourcesByNameOrGroup(ctx, input, language)
	if err != nil {
		return nil, ErrInternal
	}

	tableResponses := make([]models.SampleSourceAdminTableResponse, len(sampleSources))
	for i, sample := range sampleSources {
		tableResponses[i] = sample.ToAdminTableResponse(language)
	}
	return tableResponses, nil
}

func (s *sampleSourceService) Create(ctx context.Context, input models.SampleSourceCreateInput) (*models.SampleSourceAdminDetailResponse, error) {
	sampleSource := models.SampleSource{
		Names:    input.Names,
		Groups:   input.Groups,
		IsActive: input.IsActive,
	}

	existingSampleSource, err := s.Repo.GetSampleSourceDuplicate(ctx, sampleSource.Names, uuid.UUID{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}

	if existingSampleSource != nil {
		return nil, ErrConflict
	}

	if err := s.Repo.CreateSampleSource(ctx, &sampleSource); err != nil {
		return nil, ErrInternal
	}

	detailResponse := sampleSource.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *sampleSourceService) Update(ctx context.Context, ID uuid.UUID, input models.SampleSourceUpdateInput) (*models.SampleSourceAdminDetailResponse, error) {
	existingSampleSource, err := s.Repo.GetSampleSourceByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	validations.ApplySampleSourceUpdate(existingSampleSource, &input)

	if input.Names != nil {
		duplicate, err := s.Repo.GetSampleSourceDuplicate(ctx, input.Names, ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInternal
		}

		if duplicate != nil {
			return nil, ErrConflict
		}
	}

	if err := s.Repo.UpdateSampleSource(ctx, existingSampleSource); err != nil {
		return nil, ErrInternal
	}

	detailResponse := existingSampleSource.ToAdminDetailResponse()
	return &detailResponse, nil
}

func (s *sampleSourceService) Delete(ctx context.Context, ID uuid.UUID) error {
	sampleSource, err := s.Repo.GetSampleSourceByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return ErrInternal
	}

	if err := s.Repo.DeleteSampleSource(ctx, sampleSource); err != nil {
		return ErrInternal
	}

	return nil
}
