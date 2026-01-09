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

type mockLaboratoryRepository struct {
	GetLaboratoriesFunc                     func(ctx context.Context) ([]models.Laboratory, error)
	GetActiveLaboratoriesFunc               func(ctx context.Context) ([]models.Laboratory, error)
	GetLaboratoryByIDFunc                   func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error)
	GetLaboratoriesByNameOrAbbreviationFunc func(ctx context.Context, input string) ([]models.Laboratory, error)
	GetLaboratoryDuplicateFunc              func(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error)
	CreateLaboratoryFunc                    func(ctx context.Context, lab *models.Laboratory) error
	UpdateLaboratoryFunc                    func(ctx context.Context, lab *models.Laboratory) error
	DeleteLaboratoryFunc                    func(ctx context.Context, lab *models.Laboratory) error
}

func (r *mockLaboratoryRepository) GetLaboratories(ctx context.Context) ([]models.Laboratory, error) {
	if r.GetLaboratoriesFunc != nil {
		return r.GetLaboratoriesFunc(ctx)
	}

	return nil, nil
}

func (r *mockLaboratoryRepository) GetActiveLaboratories(ctx context.Context) ([]models.Laboratory, error) {
	if r.GetActiveLaboratoriesFunc != nil {
		return r.GetActiveLaboratoriesFunc(ctx)
	}

	return nil, nil
}

func (r *mockLaboratoryRepository) GetLaboratoryByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
	if r.GetLaboratoryByIDFunc != nil {
		return r.GetLaboratoryByIDFunc(ctx, ID)
	}

	return nil, nil
}

func (r *mockLaboratoryRepository) GetLaboratoriesByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error) {
	if r.GetLaboratoriesByNameOrAbbreviationFunc != nil {
		return r.GetLaboratoriesByNameOrAbbreviationFunc(ctx, input)
	}

	return nil, nil
}

func (r *mockLaboratoryRepository) GetLaboratoryDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
	if r.GetLaboratoryDuplicateFunc != nil {
		return r.GetLaboratoryDuplicateFunc(ctx, name, ID)
	}

	return nil, nil
}

func (r *mockLaboratoryRepository) CreateLaboratory(ctx context.Context, lab *models.Laboratory) error {
	if r.CreateLaboratoryFunc != nil {
		return r.CreateLaboratoryFunc(ctx, lab)
	}

	return nil
}

func (r *mockLaboratoryRepository) UpdateLaboratory(ctx context.Context, lab *models.Laboratory) error {
	if r.UpdateLaboratoryFunc != nil {
		return r.UpdateLaboratoryFunc(ctx, lab)
	}

	return nil
}

func (r *mockLaboratoryRepository) DeleteLaboratory(ctx context.Context, lab *models.Laboratory) error {
	if r.DeleteLaboratoryFunc != nil {
		return r.DeleteLaboratoryFunc(ctx, lab)
	}

	return nil
}

func TestLaboratoryFindAll(t *testing.T) {
	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)

	t.Run("Success", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
				return []models.Laboratory{
					lab,
				}, nil
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		expected := []models.LaboratoryAdminTableResponse{lab.ToAdminTableResponse()}

		labs, err := service.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		labs, err := service.FindAll(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, labs)
	})
}

func TestLaboratoryFindAllActive(t *testing.T) {
	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)

	t.Run("Success", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetActiveLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
				return []models.Laboratory{
					lab,
				}, nil
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		expected := []models.LaboratoryFormResponse{lab.ToFormResponse()}

		labs, err := service.FindAllActive(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetActiveLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		labs, err := service.FindAllActive(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, labs)
	})
}
func TestLaboratoryFindByID(t *testing.T) {
	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)

	t.Run("Success", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &lab, nil
			},
		}
		service := services.NewLaboratoryService(&labRepo)

		expected := lab.ToAdminTableResponse()
		labFound, err := service.FindByID(context.Background(), lab.ID)

		assert.NoError(t, err)
		assert.Equal(t, &expected, labFound)
	})

	t.Run("Error - Not found", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		lab, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, lab)
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		lab, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, lab)
	})
}

func TestLaboratoryFindByNameOrAbbreviation(t *testing.T) {
	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)

	t.Run("Success", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoriesByNameOrAbbreviationFunc: func(ctx context.Context, input string) ([]models.Laboratory, error) {
				return []models.Laboratory{lab}, nil
			},
		}
		service := services.NewLaboratoryService(&labRepo)

		expected := []models.LaboratoryAdminTableResponse{lab.ToAdminTableResponse()}
		labs, err := service.FindByNameOrAbbreviation(context.Background(), "lab")

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoriesByNameOrAbbreviationFunc: func(ctx context.Context, input string) ([]models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		labs, err := service.FindByNameOrAbbreviation(context.Background(), "lab")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, labs)
	})
}

func TestLaboratoryCreate(t *testing.T) {
	input := models.LaboratoryCreateInput{
		Name:         "Lab1",
		Abbreviation: "L1",
		IsActive:     true,
	}

	t.Run("Success", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			CreateLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return nil
			},
		}
		service := services.NewLaboratoryService(&labRepo)

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
		labRepo := mockLaboratoryRepository{
			GetLaboratoryDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoryDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
	})

	t.Run("Error - Create", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			CreateLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
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
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{ID: id}, nil
			},
			UpdateLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return nil
			},
		}
		service := services.NewLaboratoryService(&labRepo)

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
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		result, err := service.Update(context.Background(), uuid.New(), models.LaboratoryUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
			GetLaboratoryDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		result, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, result)
	})

	t.Run("Error - Update", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
			UpdateLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		result, err := service.Update(context.Background(), uuid.New(), models.LaboratoryUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
	})
}

func TestLaboratoryDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &models.Laboratory{}, nil
			},
			DeleteLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return nil
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("Error", func(t *testing.T) {
		labRepo := mockLaboratoryRepository{
			DeleteLaboratoryFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewLaboratoryService(&labRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
	})
}
