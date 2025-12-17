package sequencer_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockSequencerService struct {
	FindAllFunc            func(ctx context.Context) ([]models.Sequencer, error)
	FindAllActiveFunc      func(ctx context.Context) ([]models.SequencerFormResponse, error)
	FindByIDFunc           func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error)
	FindByBrandOrModelFunc func(ctx context.Context, input string) ([]models.Sequencer, error)
	CreateFunc             func(ctx context.Context, sequencer *models.Sequencer) error
	UpdateFunc             func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.Sequencer, error)
	DeleteFunc             func(ctx context.Context, ID uuid.UUID) error
}

func (s *MockSequencerService) FindAll(ctx context.Context) ([]models.Sequencer, error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx)
	}
	return nil, nil
}

func (s *MockSequencerService) FindAllActive(ctx context.Context) ([]models.SequencerFormResponse, error) {
	if s.FindAllActiveFunc != nil {
		return s.FindAllActiveFunc(ctx)
	}
	return nil, nil
}

func (s *MockSequencerService) FindByID(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
	if s.FindByIDFunc != nil {
		return s.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (s *MockSequencerService) FindByBrandOrModel(ctx context.Context, input string) ([]models.Sequencer, error) {
	if s.FindByBrandOrModelFunc != nil {
		return s.FindByBrandOrModelFunc(ctx, input)
	}
	return nil, nil
}

func (s *MockSequencerService) Create(ctx context.Context, sequencer *models.Sequencer) error {
	if s.CreateFunc != nil {
		return s.CreateFunc(ctx, sequencer)
	}
	return nil
}

func (s *MockSequencerService) Update(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.Sequencer, error) {
	if s.UpdateFunc != nil {
		return s.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (s *MockSequencerService) Delete(ctx context.Context, ID uuid.UUID) error {
	if s.DeleteFunc != nil {
		return s.DeleteFunc(ctx, ID)
	}
	return nil
}

func TestGetActiveSequencers(t *testing.T) {
	testutils.SetupTestContext()
	mockSequencer := models.SequencerFormResponse{ID: uuid.New()}

	t.Run("Success", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			FindAllActiveFunc: func(ctx context.Context) ([]models.SequencerFormResponse, error) {
				return []models.SequencerFormResponse{mockSequencer}, nil
			},
		}

		handler := sequencer.NewSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/sequencer", "", nil, nil,
		)
		handler.GetActiveSequencers(c)

		expected := testutils.ToJSON(
			map[string][]models.SequencerFormResponse{
				"data": {mockSequencer},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			FindAllActiveFunc: func(ctx context.Context) ([]models.SequencerFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := sequencer.NewSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/laboratory", "",
			nil, nil,
		)
		handler.GetActiveSequencers(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
