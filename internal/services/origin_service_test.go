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

func TestOriginFindAll(t *testing.T) {
	origin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Humano", "en": "Human"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
				return []models.Origin{origin}, nil
			},
		}

		service := services.NewOriginService(originRepo, nil)
		result, err := service.FindAll(context.Background(), "pt")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, origin.ToAdminTableResponse("pt"), result[0])
	})

	t.Run("Error", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.FindAll(context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestOriginFindAllActive(t *testing.T) {
	origin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Humano"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetActiveOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
				return []models.Origin{origin}, nil
			},
		}

		service := services.NewOriginService(originRepo, nil)
		result, err := service.FindAllActive(context.Background(), "pt")

		assert.NoError(t, err)
		assert.Equal(t, []models.OriginFormResponse{
			origin.ToFormResponse("pt"),
		}, result)
	})

	t.Run("Error", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetActiveOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.FindAllActive(context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestOriginFindByID(t *testing.T) {
	origin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Humano"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &origin, nil
			},
		}

		service := services.NewOriginService(originRepo, nil)
		result, err := service.FindByID(context.Background(), origin.ID)

		assert.NoError(t, err)
		assert.Equal(t, origin.ToAdminDetailResponse(), *result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestOriginFindByName(t *testing.T) {
	origin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"en": "Human"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginsByNameFunc: func(ctx context.Context, name, lang string) ([]models.Origin, error) {
				return []models.Origin{origin}, nil
			},
		}

		service := services.NewOriginService(originRepo, nil)
		result, err := service.FindByName(context.Background(), "hum", "en")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, origin.ToAdminTableResponse("en"), result[0])
	})

	t.Run("Error", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginsByNameFunc: func(ctx context.Context, name, lang string) ([]models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.FindByName(context.Background(), "hum", "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestOriginCreate(t *testing.T) {
	input := models.OriginCreateInput{
		Names:    models.JSONMap{"pt": "Humano"},
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return nil
			},
		}

		service := services.NewOriginService(originRepo, nil)
		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - Find duplicate", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Create", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestOriginUpdate(t *testing.T) {
	id := uuid.New()

	isActive := true
	input := models.OriginUpdateInput{
		Names:    map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		IsActive: &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{ID: id}, nil
			},
			UpdateOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return nil
			},
		}

		service := services.NewOriginService(originRepo, nil)
		result, err := service.Update(context.Background(), id, models.OriginUpdateInput{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.Update(context.Background(), id, models.OriginUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{ID: id}, nil
			},
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{ID: id}, nil
			},
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
			UpdateOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		result, err := service.Update(context.Background(), id, models.OriginUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestOriginDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
			DeleteOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return nil
			},
		}

		service := services.NewOriginService(originRepo, nil)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Delete", func(t *testing.T) {
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
			DeleteOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewOriginService(originRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
