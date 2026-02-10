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

func TestSampleSourceFindAll(t *testing.T) {
	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true,
	)
	language := "en"

	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
				return []models.SampleSource{sampleSource}, nil
			},
		}

		service := services.NewSampleSourceService(sampleSourceRepo, nil)
		expected := []models.SampleSourceAdminTableResponse{
			sampleSource.ToAdminTableResponse(language),
		}

		sampleSources, err := service.FindAll(context.Background(), language)

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		sampleSources, err := service.FindAll(context.Background(), language)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSources)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleSourceFindAllActive(t *testing.T) {
	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetActiveSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
				return []models.SampleSource{sampleSource}, nil
			},
		}

		service := services.NewSampleSourceService(sampleSourceRepo, nil)
		expected := []models.SampleSourceFormResponse{sampleSource.ToFormResponse("en")}

		sampleSources, err := service.FindAllActive(context.Background(), "en")

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetActiveSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		sampleSources, err := service.FindAllActive(context.Background(), "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSources)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleSourceFindByID(t *testing.T) {
	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &sampleSource, nil
			},
		}

		expected := sampleSource.ToAdminDetailResponse()
		service := services.NewSampleSourceService(sampleSourceRepo, nil)
		sampleSourceFound, err := service.FindByID(context.Background(), sampleSource.ID)

		assert.NoError(t, err)
		assert.Equal(t, &expected, sampleSourceFound)
	})

	t.Run("Record not found", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		sampleSourceFound, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, sampleSourceFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("DB error", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		sampleSourceFound, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSourceFound)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleSourceFindByNameOrGroup(t *testing.T) {
	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true,
	)
	language := "en"

	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourcesByNameOrGroupFunc: func(ctx context.Context, input, language string) ([]models.SampleSource, error) {
				return []models.SampleSource{sampleSource}, nil
			},
		}

		expected := []models.SampleSourceAdminTableResponse{
			sampleSource.ToAdminTableResponse(language),
		}
		service := services.NewSampleSourceService(sampleSourceRepo, nil)
		sampleSources, err := service.FindByNameOrGroup(context.Background(), "plasma", language)

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourcesByNameOrGroupFunc: func(ctx context.Context, input, language string) ([]models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		sampleSources, err := service.FindByNameOrGroup(context.Background(), "plasma", "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSources)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleSourceCreate(t *testing.T) {
	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			CreateSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return nil
			},
		}

		expected := sampleSource.ToAdminDetailResponse()
		service := services.NewSampleSourceService(sampleSourceRepo, nil)
		result, err := service.Create(
			context.Background(),
			models.SampleSourceCreateInput{
				Names:    sampleSource.Names,
				Groups:   sampleSource.Groups,
				IsActive: sampleSource.IsActive,
			},
		)
		expected.ID = uuid.Nil

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Find duplicate", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		result, err := service.Create(
			context.Background(),
			models.SampleSourceCreateInput{
				Names:    sampleSource.Names,
				Groups:   sampleSource.Groups,
				IsActive: sampleSource.IsActive,
			},
		)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		result, err := service.Create(
			context.Background(),
			models.SampleSourceCreateInput{
				Names:    sampleSource.Names,
				Groups:   sampleSource.Groups,
				IsActive: sampleSource.IsActive,
			},
		)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Create", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			CreateSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		result, err := service.Create(
			context.Background(),
			models.SampleSourceCreateInput{
				Names:    sampleSource.Names,
				Groups:   sampleSource.Groups,
				IsActive: sampleSource.IsActive,
			},
		)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleSourceUpdate(t *testing.T) {
	id := uuid.New()
	isActive := true

	input := models.SampleSourceUpdateInput{
		Names:    map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		Groups:   map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		IsActive: &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{ID: uuid.New()}, nil
			},
			UpdateSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return nil
			},
		}

		service := services.NewSampleSourceService(sampleSourceRepo, nil)
		sampleSource, err := service.Update(context.Background(), uuid.New(), models.SampleSourceUpdateInput{})

		assert.NoError(t, err)
		assert.NotEmpty(t, sampleSource)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		sampleSource, err := service.Update(context.Background(), uuid.New(), models.SampleSourceUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, sampleSource)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
			GetSampleSourceDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		sampleSource, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, sampleSource)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
			UpdateSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		sampleSource, err := service.Update(context.Background(), uuid.New(), models.SampleSourceUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSource)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleSourceDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
			DeleteSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return nil
			},
		}

		service := services.NewSampleSourceService(sampleSourceRepo, nil)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})
	t.Run("Error - Not Found", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Delete", func(t *testing.T) {
		sampleSourceRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
			DeleteSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewSampleSourceService(sampleSourceRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
