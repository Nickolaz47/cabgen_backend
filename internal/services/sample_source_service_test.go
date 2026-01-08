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

type mockSampleSourceRepository struct {
	GetSampleSourcesFunc              func(ctx context.Context) ([]models.SampleSource, error)
	GetActiveSampleSourcesFunc        func(ctx context.Context) ([]models.SampleSource, error)
	GetSampleSourceByIDFunc           func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error)
	GetSampleSourcesByNameOrGroupFunc func(ctx context.Context, input, language string) ([]models.SampleSource, error)
	GetSampleSourceDuplicateFunc      func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error)
	CreateSampleSourceFunc            func(ctx context.Context, sampleSource *models.SampleSource) error
	UpdateSampleSourceFunc            func(ctx context.Context, sampleSource *models.SampleSource) error
	DeleteSampleSourceFunc            func(ctx context.Context, sampleSource *models.SampleSource) error
}

func (r *mockSampleSourceRepository) GetSampleSources(ctx context.Context) ([]models.SampleSource, error) {
	if r.GetSampleSourcesFunc != nil {
		return r.GetSampleSourcesFunc(ctx)
	}

	return nil, nil
}

func (r *mockSampleSourceRepository) GetActiveSampleSources(ctx context.Context) ([]models.SampleSource, error) {
	if r.GetActiveSampleSourcesFunc != nil {
		return r.GetActiveSampleSourcesFunc(ctx)
	}

	return nil, nil
}

func (r *mockSampleSourceRepository) GetSampleSourceByID(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
	if r.GetSampleSourceByIDFunc != nil {
		return r.GetSampleSourceByIDFunc(ctx, ID)
	}

	return nil, nil
}

func (r *mockSampleSourceRepository) GetSampleSourcesByNameOrGroup(ctx context.Context, input, language string) ([]models.SampleSource, error) {
	if r.GetSampleSourcesByNameOrGroupFunc != nil {
		return r.GetSampleSourcesByNameOrGroupFunc(ctx, input, language)
	}

	return nil, nil
}

func (r *mockSampleSourceRepository) GetSampleSourceDuplicate(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
	if r.GetSampleSourceDuplicateFunc != nil {
		return r.GetSampleSourceDuplicateFunc(ctx, names, ID)
	}

	return nil, nil
}

func (r *mockSampleSourceRepository) CreateSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	if r.CreateSampleSourceFunc != nil {
		return r.CreateSampleSourceFunc(ctx, sampleSource)
	}

	return nil
}

func (r *mockSampleSourceRepository) UpdateSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	if r.UpdateSampleSourceFunc != nil {
		return r.UpdateSampleSourceFunc(ctx, sampleSource)
	}

	return nil
}

func (r *mockSampleSourceRepository) DeleteSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	if r.DeleteSampleSourceFunc != nil {
		return r.DeleteSampleSourceFunc(ctx, sampleSource)
	}

	return nil
}

func TestSampleSourceFindAll(t *testing.T) {
	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true,
	)
	language := "en"

	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
				return []models.SampleSource{sampleSource}, nil
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		expected := []models.SampleSourceAdminTableResponse{
			sampleSource.ToAdminTableResponse(language),
		}

		sampleSources, err := service.FindAll(context.Background(), language)

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSources, err := service.FindAll(context.Background(), language)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSources)
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
		sampleSourceRepo := mockSampleSourceRepository{
			GetActiveSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
				return []models.SampleSource{sampleSource}, nil
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		expected := []models.SampleSourceFormResponse{sampleSource.ToFormResponse("en")}

		sampleSources, err := service.FindAllActive(context.Background(), "en")

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetActiveSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSources, err := service.FindAllActive(context.Background(), "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSources)
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
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &sampleSource, nil
			},
		}

		expected := sampleSource.ToAdminDetailResponse()
		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSourceFound, err := service.FindByID(context.Background(), sampleSource.ID)

		assert.NoError(t, err)
		assert.Equal(t, &expected, sampleSourceFound)
	})

	t.Run("Record not found", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSourceFound, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, sampleSourceFound)
	})

	t.Run("DB error", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSourceFound, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSourceFound)
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
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourcesByNameOrGroupFunc: func(ctx context.Context, input, language string) ([]models.SampleSource, error) {
				return []models.SampleSource{sampleSource}, nil
			},
		}

		expected := []models.SampleSourceAdminTableResponse{
			sampleSource.ToAdminTableResponse(language),
		}
		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSources, err := service.FindByNameOrGroup(context.Background(), "plasma", language)

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourcesByNameOrGroupFunc: func(ctx context.Context, input, language string) ([]models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSources, err := service.FindByNameOrGroup(context.Background(), "plasma", "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSources)
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
		sampleSourceRepo := mockSampleSourceRepository{
			CreateSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return nil
			},
		}

		expected := sampleSource.ToAdminDetailResponse()
		service := services.NewSampleSourceService(&sampleSourceRepo)
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
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
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
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
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
	})

	t.Run("Error - Create", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			CreateSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
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
	})
}

func TestSampleSourceUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{ID: uuid.New()}, nil
			},
			UpdateSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return nil
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSource, err := service.Update(context.Background(), uuid.New(), models.SampleSourceUpdateInput{})

		assert.NoError(t, err)
		assert.NotEmpty(t, sampleSource)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSource, err := service.Update(context.Background(), uuid.New(), models.SampleSourceUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, sampleSource)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
			GetSampleSourceDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSource, err := service.Update(context.Background(), uuid.New(), models.SampleSourceUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, sampleSource)
	})

	t.Run("Error - Update", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
			UpdateSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		sampleSource, err := service.Update(context.Background(), uuid.New(), models.SampleSourceUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, sampleSource)
	})
}

func TestSampleSourceDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
			DeleteSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return nil
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})
	t.Run("Error - Not Found", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("Error - Delete", func(t *testing.T) {
		sampleSourceRepo := mockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
				return &models.SampleSource{}, nil
			},
			DeleteSampleSourceFunc: func(ctx context.Context, sampleSource *models.SampleSource) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewSampleSourceService(&sampleSourceRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
	})
}
