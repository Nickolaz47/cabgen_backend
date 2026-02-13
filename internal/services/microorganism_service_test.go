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

func TestMicroorganismFindAll(t *testing.T) {
	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B"},
		true,
	)
	language := "en"

	t.Run("Success", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismsFunc: func(ctx context.Context) ([]models.Microorganism, error) {
				return []models.Microorganism{micro}, nil
			},
		}

		service := services.NewMicroorganismService(microRepo, nil)
		expected := []models.MicroorganismAdminTableResponse{
			micro.ToAdminTableResponse(language),
		}

		micros, err := service.FindAll(context.Background(), language)

		assert.NoError(t, err)
		assert.Equal(t, expected, micros)
	})

	t.Run("Error", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismsFunc: func(ctx context.Context) ([]models.Microorganism, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		micros, err := service.FindAll(context.Background(), language)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, micros)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestMicroorganismFindAllActive(t *testing.T) {
	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetActiveMicroorganismsFunc: func(ctx context.Context) ([]models.Microorganism, error) {
				return []models.Microorganism{micro}, nil
			},
		}

		service := services.NewMicroorganismService(microRepo, nil)
		expected := []models.MicroorganismFormResponse{micro.ToFormResponse("en")}

		micros, err := service.FindAllActive(context.Background(), "en")

		assert.NoError(t, err)
		assert.Equal(t, expected, micros)
	})

	t.Run("Error", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetActiveMicroorganismsFunc: func(ctx context.Context) ([]models.Microorganism, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		micros, err := service.FindAllActive(context.Background(), "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, micros)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestMicroorganismFindByID(t *testing.T) {
	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return &micro, nil
			},
		}

		expected := micro.ToAdminDetailResponse()
		service := services.NewMicroorganismService(microRepo, nil)
		microFound, err := service.FindByID(context.Background(), micro.ID)

		assert.NoError(t, err)
		assert.Equal(t, &expected, microFound)
	})

	t.Run("Record not found", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		microFound, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, microFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("DB error", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		microFound, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, microFound)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestMicroorganismFindBySpecies(t *testing.T) {
	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B"},
		true,
	)
	language := "en"

	t.Run("Success", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismsBySpeciesFunc: func(ctx context.Context, input, language string) ([]models.Microorganism, error) {
				return []models.Microorganism{micro}, nil
			},
		}

		expected := []models.MicroorganismAdminTableResponse{
			micro.ToAdminTableResponse(language),
		}
		service := services.NewMicroorganismService(microRepo, nil)
		micros, err := service.FindBySpecies(context.Background(), "meningitidis", language)

		assert.NoError(t, err)
		assert.Equal(t, expected, micros)
	})

	t.Run("Error", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismsBySpeciesFunc: func(ctx context.Context, input, language string) ([]models.Microorganism, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		micros, err := service.FindBySpecies(context.Background(), "meningitidis", "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, micros)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestMicroorganismCreate(t *testing.T) {
	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismDuplicateFunc: func(ctx context.Context, species string, variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateMicroorganismFunc: func(ctx context.Context, micro *models.Microorganism) error {
				return nil
			},
		}

		expected := micro.ToAdminDetailResponse()
		service := services.NewMicroorganismService(microRepo, nil)
		result, err := service.Create(
			context.Background(),
			models.MicroorganismCreateInput{
				Taxon:    micro.Taxon,
				Species:  micro.Species,
				Variety:  micro.Variety,
				IsActive: micro.IsActive,
			},
		)
		expected.ID = uuid.Nil

		assert.NoError(t, err)
		// Ajuste para ignorar ID gerado aleatoriamente se necessário, ou assumir que o ToAdminDetailResponse lida com isso
		// No exemplo original, ele compara os valores retornados.
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Find duplicate", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismDuplicateFunc: func(ctx context.Context, species string, variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		result, err := service.Create(
			context.Background(),
			models.MicroorganismCreateInput{
				Taxon:    micro.Taxon,
				Species:  micro.Species,
				Variety:  micro.Variety,
				IsActive: micro.IsActive,
			},
		)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismDuplicateFunc: func(ctx context.Context, species string, variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error) {
				return &models.Microorganism{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		result, err := service.Create(
			context.Background(),
			models.MicroorganismCreateInput{
				Taxon:    micro.Taxon,
				Species:  micro.Species,
				Variety:  micro.Variety,
				IsActive: micro.IsActive,
			},
		)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Create", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismDuplicateFunc: func(ctx context.Context, species string, variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateMicroorganismFunc: func(ctx context.Context, micro *models.Microorganism) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		result, err := service.Create(
			context.Background(),
			models.MicroorganismCreateInput{
				Taxon:    micro.Taxon,
				Species:  micro.Species,
				Variety:  micro.Variety,
				IsActive: micro.IsActive,
			},
		)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestMicroorganismUpdate(t *testing.T) {
	id := uuid.New()
	isActive := true
	species := "Neisseria Updated"

	input := models.MicroorganismUpdateInput{
		Species:  &species,
		Variety:  map[string]string{"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B"},
		IsActive: &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return &models.Microorganism{ID: uuid.New()}, nil
			},
			GetMicroorganismDuplicateFunc: func(ctx context.Context, species string, variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrRecordNotFound
			},
			UpdateMicroorganismFunc: func(ctx context.Context, micro *models.Microorganism) error {
				return nil
			},
		}

		service := services.NewMicroorganismService(microRepo, nil)
		micro, err := service.Update(context.Background(), uuid.New(), models.MicroorganismUpdateInput{})

		assert.NoError(t, err)
		assert.NotEmpty(t, micro)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		micro, err := service.Update(context.Background(), uuid.New(), models.MicroorganismUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, micro)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return &models.Microorganism{}, nil
			},
			GetMicroorganismDuplicateFunc: func(ctx context.Context, species string, variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error) {
				return &models.Microorganism{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		micro, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, micro)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return &models.Microorganism{}, nil
			},
			UpdateMicroorganismFunc: func(ctx context.Context, micro *models.Microorganism) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		micro, err := service.Update(context.Background(), uuid.New(), models.MicroorganismUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, micro)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestMicroorganismDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return &models.Microorganism{}, nil
			},
			DeleteMicroorganismFunc: func(ctx context.Context, micro *models.Microorganism) error {
				return nil
			},
		}

		service := services.NewMicroorganismService(microRepo, nil)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Delete", func(t *testing.T) {
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
				return &models.Microorganism{}, nil
			},
			DeleteMicroorganismFunc: func(ctx context.Context, micro *models.Microorganism) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewMicroorganismService(microRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
