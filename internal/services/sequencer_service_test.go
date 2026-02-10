package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestSequencerFindAll(t *testing.T) {
	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)

	t.Run("Success", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return []models.Sequencer{sequencer}, nil
			},
		}
		service := services.NewSequencerService(seqRepo, nil)

		expected := []models.SequencerAdminTableResponse{sequencer.ToAdminTableResponse()}
		sequencers, err := service.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, sequencers)
	})

	t.Run("Error", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		sequencers, err := service.FindAll(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sequencers)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSequencerFindAllActive(t *testing.T) {
	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)

	t.Run("Success", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetActiveSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return []models.Sequencer{sequencer}, nil
			},
		}
		service := services.NewSequencerService(seqRepo, nil)

		expected := []models.SequencerFormResponse{sequencer.ToFormResponse()}
		sequencers, err := service.FindAllActive(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, sequencers)
	})

	t.Run("Error", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetActiveSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		sequencers, err := service.FindAllActive(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sequencers)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSequencerFindByID(t *testing.T) {
	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)

	t.Run("Success", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &sequencer, nil
			},
		}
		service := services.NewSequencerService(seqRepo, nil)

		expected := sequencer.ToAdminTableResponse()
		result, err := service.FindByID(context.Background(), sequencer.ID)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Record not found", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.FindByID(context.Background(), sequencer.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSequencerFindByBrandOrModel(t *testing.T) {
	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)

	t.Run("Success", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencersByBrandOrModelFunc: func(ctx context.Context, input string) ([]models.Sequencer, error) {
				return []models.Sequencer{sequencer}, nil
			},
		}
		service := services.NewSequencerService(seqRepo, nil)

		expected := []models.SequencerAdminTableResponse{sequencer.ToAdminTableResponse()}
		result, err := service.FindByBrandOrModel(context.Background(), "illumin")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencersByBrandOrModelFunc: func(ctx context.Context, input string) ([]models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.FindByBrandOrModel(context.Background(), "illumin")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSequencerCreate(t *testing.T) {
	input := models.SequencerCreateInput{
		Model:    "MiSeq",
		Brand:    "Illumina",
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			CreateSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return nil
			},
		}
		service := services.NewSequencerService(seqRepo, nil)

		expected := models.SequencerAdminTableResponse{
			Model:    input.Model,
			Brand:    input.Brand,
			IsActive: input.IsActive,
		}
		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Find duplicate", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerDuplicateFunc: func(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerDuplicateFunc: func(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Create", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			CreateSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSequencerUpdate(t *testing.T) {
	id := uuid.New()
	model, brand, isActive := "MiSeq", "Illumina", true
	input := models.SequencerUpdateInput{
		Model:    &model,
		Brand:    &brand,
		IsActive: &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{ID: id}, nil
			},
			UpdateSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return nil
			},
		}
		service := services.NewSequencerService(seqRepo, nil)

		expected := models.SequencerAdminTableResponse{
			ID:       id,
			Model:    *input.Model,
			Brand:    *input.Brand,
			IsActive: *input.IsActive,
		}
		result, err := service.Update(context.Background(), id, input)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.Update(context.Background(), uuid.New(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{ID: uuid.New()}, nil
			},
			GetSequencerDuplicateFunc: func(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{ID: uuid.New()}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.Update(context.Background(), uuid.New(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{ID: uuid.New()}, nil
			},
			UpdateSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		result, err := service.Update(context.Background(), uuid.New(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSequencerDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{}, nil
			},
			DeleteSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return nil
			},
		}

		service := services.NewSequencerService(seqRepo, nil)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error", func(t *testing.T) {
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{}, nil
			},
			DeleteSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSequencerService(seqRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
