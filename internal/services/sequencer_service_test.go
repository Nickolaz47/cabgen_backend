package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockSequencerRepository struct {
	GetSequencersFunc               func(ctx context.Context) ([]models.Sequencer, error)
	GetActiveSequencersFunc         func(ctx context.Context) ([]models.Sequencer, error)
	GetSequencerByIDFunc            func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error)
	GetSequencersByBrandOrModelFunc func(ctx context.Context, input string) ([]models.Sequencer, error)
	GetSequencerDuplicateFunc       func(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error)
	CreateSequencerFunc             func(ctx context.Context, sequencer *models.Sequencer) error
	UpdateSequencerFunc             func(ctx context.Context, sequencer *models.Sequencer) error
	DeleteSequencerFunc             func(ctx context.Context, sequencer *models.Sequencer) error
}

func (s *mockSequencerRepository) GetSequencers(ctx context.Context) ([]models.Sequencer, error) {
	if s.GetSequencersFunc != nil {
		return s.GetSequencersFunc(ctx)
	}
	return nil, nil
}

func (s *mockSequencerRepository) GetActiveSequencers(ctx context.Context) ([]models.Sequencer, error) {
	if s.GetActiveSequencersFunc != nil {
		return s.GetActiveSequencersFunc(ctx)
	}
	return nil, nil
}

func (s *mockSequencerRepository) GetSequencerByID(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
	if s.GetSequencerByIDFunc != nil {
		return s.GetSequencerByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (s *mockSequencerRepository) GetSequencersByBrandOrModel(ctx context.Context, input string) ([]models.Sequencer, error) {
	if s.GetSequencersByBrandOrModelFunc != nil {
		return s.GetSequencersByBrandOrModelFunc(ctx, input)
	}
	return nil, nil
}

func (s *mockSequencerRepository) GetSequencerDuplicate(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
	if s.GetSequencerDuplicateFunc != nil {
		return s.GetSequencerDuplicateFunc(ctx, model, ID)
	}
	return nil, nil
}

func (s *mockSequencerRepository) CreateSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	if s.CreateSequencerFunc != nil {
		return s.CreateSequencerFunc(ctx, sequencer)
	}
	return nil
}

func (s *mockSequencerRepository) UpdateSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	if s.UpdateSequencerFunc != nil {
		return s.UpdateSequencerFunc(ctx, sequencer)
	}
	return nil
}

func (s *mockSequencerRepository) DeleteSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	if s.DeleteSequencerFunc != nil {
		return s.DeleteSequencerFunc(ctx, sequencer)
	}
	return nil
}

func TestSequencerFindAll(t *testing.T) {
	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)

	t.Run("Success", func(t *testing.T) {
		seqRepo := mockSequencerRepository{
			GetSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return []models.Sequencer{sequencer}, nil
			},
		}
		service := services.NewSequencerService(&seqRepo)

		expected := []models.SequencerAdminTableResponse{sequencer.ToAdminTableResponse()}
		sequencers, err := service.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, sequencers)
	})

	t.Run("Error", func(t *testing.T) {
		seqRepo := mockSequencerRepository{
			GetSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		service := services.NewSequencerService(&seqRepo)

		sequencers, err := service.FindAll(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sequencers)
	})
}

func TestSequencerFindAllActive(t *testing.T) {
	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)

	t.Run("Success", func(t *testing.T) {
		seqRepo := mockSequencerRepository{
			GetActiveSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return []models.Sequencer{sequencer}, nil
			},
		}
		service := services.NewSequencerService(&seqRepo)

		expected := []models.SequencerFormResponse{sequencer.ToFormResponse()}
		sequencers, err := service.FindAllActive(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, sequencers)
	})

	t.Run("Error", func(t *testing.T) {
		seqRepo := mockSequencerRepository{
			GetActiveSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		service := services.NewSequencerService(&seqRepo)

		sequencers, err := service.FindAllActive(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sequencers)
	})
}

func TestSequencerFindByID(t *testing.T) {
	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)

	t.Run("Success", func(t *testing.T) {
		seqRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &sequencer, nil
			},
		}
		service := services.NewSequencerService(&seqRepo)

		expected := sequencer.ToAdminTableResponse()
		result, err := service.FindByID(context.Background(), sequencer.ID)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Record not found", func(t *testing.T) {
		seqRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewSequencerService(&seqRepo)
		result, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		seqRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		service := services.NewSequencerService(&seqRepo)

		result, err := service.FindByID(context.Background(), sequencer.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
	})
}

func TestSequencerFindByBrandOrModel(t *testing.T) {
	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)

	t.Run("Success", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			GetSequencersByBrandOrModelFunc: func(ctx context.Context, input string) ([]models.Sequencer, error) {
				return []models.Sequencer{sequencer}, nil
			},
		}
		service := services.NewSequencerService(&sequencerRepo)

		expected := []models.SequencerAdminTableResponse{sequencer.ToAdminTableResponse()}
		result, err := service.FindByBrandOrModel(context.Background(), "illumin")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			GetSequencersByBrandOrModelFunc: func(ctx context.Context, input string) ([]models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		result, err := service.FindByBrandOrModel(context.Background(), "illumin")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
	})
}

func TestSequencerCreate(t *testing.T) {
	input := models.SequencerCreateInput{
		Model:    "MiSeq",
		Brand:    "Illumina",
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			CreateSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return nil
			},
		}
		service := services.NewSequencerService(&sequencerRepo)

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
		sequencerRepo := mockSequencerRepository{
			GetSequencerDuplicateFunc: func(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			GetSequencerDuplicateFunc: func(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{}, nil
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
	})

	t.Run("Error - Create", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			CreateSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
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
		sequencerRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{ID: id}, nil
			},
			UpdateSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return nil
			},
		}
		service := services.NewSequencerService(&sequencerRepo)

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
		sequencerRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		result, err := service.Update(context.Background(), uuid.New(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{ID: uuid.New()}, nil
			},
			GetSequencerDuplicateFunc: func(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{ID: uuid.New()}, nil
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		result, err := service.Update(context.Background(), uuid.New(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
	})

	t.Run("Error - Update", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{ID: uuid.New()}, nil
			},
			UpdateSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		result, err := service.Update(context.Background(), uuid.New(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
	})
}

func TestSequencerDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{}, nil
			},
			DeleteSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return nil
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("Error", func(t *testing.T) {
		sequencerRepo := mockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
				return &models.Sequencer{}, nil
			},
			DeleteSequencerFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSequencerService(&sequencerRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
	})
}
