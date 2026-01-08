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

type SequencerService interface {
	FindAll(ctx context.Context) ([]models.SequencerAdminTableResponse, error)
	FindAllActive(ctx context.Context) ([]models.SequencerFormResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*models.SequencerAdminTableResponse, error)
	FindByBrandOrModel(ctx context.Context, input string) ([]models.SequencerAdminTableResponse, error)
	Create(ctx context.Context, input models.SequencerCreateInput) (*models.SequencerAdminTableResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.SequencerAdminTableResponse, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type sequencerService struct {
	Repo repository.SequencerRepository
}

func NewSequencerService(repo repository.SequencerRepository) SequencerService {
	return &sequencerService{Repo: repo}
}

func (s *sequencerService) FindAll(ctx context.Context) ([]models.SequencerAdminTableResponse, error) {
	sequencers, err := s.Repo.GetSequencers(ctx)

	if err != nil {
		return nil, ErrInternal
	}

	tableResponses := make([]models.SequencerAdminTableResponse, len(sequencers))
	for i, sequencer := range sequencers {
		tableResponses[i] = sequencer.ToAdminTableResponse()
	}

	return tableResponses, nil
}

func (s *sequencerService) FindAllActive(ctx context.Context) ([]models.SequencerFormResponse, error) {
	sequencers, err := s.Repo.GetActiveSequencers(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	formSequencers := make([]models.SequencerFormResponse, len(sequencers))
	for i, sequencer := range sequencers {
		formSequencers[i] = sequencer.ToFormResponse()
	}

	return formSequencers, nil
}

func (s *sequencerService) FindByID(
	ctx context.Context,
	ID uuid.UUID) (*models.SequencerAdminTableResponse, error) {
	sequencer, err := s.Repo.GetSequencerByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	tableResponse := sequencer.ToAdminTableResponse()
	return &tableResponse, nil
}

func (s *sequencerService) FindByBrandOrModel(
	ctx context.Context,
	input string) ([]models.SequencerAdminTableResponse, error) {
	sequencers, err := s.Repo.GetSequencersByBrandOrModel(ctx, input)
	if err != nil {
		return nil, ErrInternal
	}

	tableResponses := make([]models.SequencerAdminTableResponse, len(sequencers))
	for i, sequencer := range sequencers {
		tableResponses[i] = sequencer.ToAdminTableResponse()
	}

	return tableResponses, nil
}

func (s *sequencerService) Create(
	ctx context.Context,
	input models.SequencerCreateInput) (*models.SequencerAdminTableResponse, error) {
	sequencer := models.Sequencer{
		Model:    input.Model,
		Brand:    input.Brand,
		IsActive: input.IsActive,
	}

	existingSequencer, err := s.Repo.GetSequencerDuplicate(ctx, sequencer.Model, uuid.UUID{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}

	if existingSequencer != nil {
		return nil, ErrConflict
	}

	if err := s.Repo.CreateSequencer(ctx, &sequencer); err != nil {
		return nil, ErrInternal
	}

	tableResponse := sequencer.ToAdminTableResponse()
	return &tableResponse, nil
}

func (s *sequencerService) Update(
	ctx context.Context,
	ID uuid.UUID,
	input models.SequencerUpdateInput) (*models.SequencerAdminTableResponse, error) {
	existingSequencer, err := s.Repo.GetSequencerByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, ErrInternal
	}

	validations.ApplySequencerUpdate(existingSequencer, &input)

	if input.Model != nil {
		duplicate, err := s.Repo.GetSequencerDuplicate(ctx, *input.Model, ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInternal
		}

		if duplicate != nil {
			return nil, ErrConflict
		}
	}

	if err := s.Repo.UpdateSequencer(ctx, existingSequencer); err != nil {
		return nil, ErrInternal
	}

	tableResponse := existingSequencer.ToAdminTableResponse()
	return &tableResponse, nil
}

func (s *sequencerService) Delete(ctx context.Context, ID uuid.UUID) error {
	sequencer, err := s.Repo.GetSequencerByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return ErrInternal
	}

	if err := s.Repo.DeleteSequencer(ctx, sequencer); err != nil {
		return ErrInternal
	}

	return nil
}
