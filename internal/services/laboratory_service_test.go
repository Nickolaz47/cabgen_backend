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

func TestLaboratoryFindAll(t *testing.T) {
	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)

	t.Run("Success", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
				return []models.Laboratory{
					lab,
				}, nil
			},
		}

		service := services.NewLaboratoryService(labRepo, nil)
		expected := []models.LaboratoryAdminTableResponse{lab.ToAdminTableResponse()}

		labs, err := service.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		labs, err := service.FindAll(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, labs)
		assert.Equal(t, logs.Len(), 1)
	})
}

func TestLaboratoryFindAllActive(t *testing.T) {
	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)

	t.Run("Success", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetActiveLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
				return []models.Laboratory{
					lab,
				}, nil
			},
		}

		service := services.NewLaboratoryService(labRepo, nil)
		expected := []models.LaboratoryFormResponse{lab.ToFormResponse()}

		labs, err := service.FindAllActive(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetActiveLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		labs, err := service.FindAllActive(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, labs)
		assert.Equal(t, logs.Len(), 1)
	})
}

func TestLaboratoryFindByID(t *testing.T) {
	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)

	t.Run("Success", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &lab, nil
			},
		}
		service := services.NewLaboratoryService(labRepo, nil)

		expected := lab.ToAdminTableResponse()
		labFound, err := service.FindByID(context.Background(), lab.ID)

		assert.NoError(t, err)
		assert.Equal(t, &expected, labFound)
	})

	t.Run("Error - Not found", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		lab, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, lab)
		assert.Equal(t, logs.Len(), 1)
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		lab, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, lab)
		assert.Equal(t, logs.Len(), 1)
	})
}

func TestLaboratoryFindByNameOrAbbreviation(t *testing.T) {
	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)

	t.Run("Success", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoriesByNameOrAbbreviationFunc: func(ctx context.Context, input string) ([]models.Laboratory, error) {
				return []models.Laboratory{lab}, nil
			},
		}
		service := services.NewLaboratoryService(labRepo, nil)

		expected := []models.LaboratoryAdminTableResponse{lab.ToAdminTableResponse()}
		labs, err := service.FindByNameOrAbbreviation(context.Background(), "lab")

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoriesByNameOrAbbreviationFunc: func(ctx context.Context, input string) ([]models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		labs, err := service.FindByNameOrAbbreviation(context.Background(), "lab")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, labs)
		assert.Equal(t, logs.Len(), 1)
	})
}

func TestLaboratoryCreate(t *testing.T) {
	input := models.LaboratoryCreateInput{
		Name:         "Lab1",
		Abbreviation: "L1",
		IsActive:     true,
	}

	t.Run("Success", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			CreateLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return nil
			},
		}
		service := services.NewLaboratoryService(labRepo, nil)

		expected := models.LaboratoryAdminTableResponse{
			Name:         input.Name,
			Abbreviation: input.Abbreviation,
			IsActive:     input.IsActive,
		}
		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Find duplicate", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, logs.Len(), 1)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
		assert.Equal(t, logs.Len(), 1)
	})

	t.Run("Error - Create", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			CreateLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, logs.Len(), 1)
	})
}

func TestLaboratoryUpdate(t *testing.T) {
	id := uuid.New()
	name, abbreviation, isActive := "Lab1", "L1", true
	input := models.LaboratoryUpdateInput{
		Name:         &name,
		Abbreviation: &abbreviation,
		IsActive:     &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{ID: id}, nil
			},
			UpdateLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return nil
			},
		}
		service := services.NewLaboratoryService(labRepo, nil)

		expected := models.LaboratoryAdminTableResponse{
			ID:           id,
			Name:         *input.Name,
			Abbreviation: *input.Abbreviation,
			IsActive:     *input.IsActive,
		}
		result, err := service.Update(context.Background(), id, input)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		result, err := service.Update(context.Background(), uuid.New(), models.LaboratoryUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, logs.Len(), 1)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
			GetLaboratoryDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		result, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
		assert.Equal(t, logs.Len(), 1)
	})

	t.Run("Error - Update", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
			UpdateLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		result, err := service.Update(context.Background(), uuid.New(), models.LaboratoryUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, logs.Len(), 1)
	})
}

func TestLaboratoryDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
			DeleteLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return nil
			},
		}

		service := services.NewLaboratoryService(labRepo, nil)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, logs.Len(), 1)
	})

	t.Run("Error", func(t *testing.T) {
		labRepo := &mocks.MockLaboratoryRepository{
			DeleteLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewLaboratoryService(labRepo, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, logs.Len(), 1)
	})
}
