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

type mockOriginRepository struct {
	GetOriginsFunc         func(ctx context.Context) ([]models.Origin, error)
	GetActiveOriginsFunc   func(ctx context.Context) ([]models.Origin, error)
	GetOriginByIDFunc      func(ctx context.Context, ID uuid.UUID) (*models.Origin, error)
	GetOriginsByNameFunc   func(ctx context.Context, name, lang string) ([]models.Origin, error)
	GetOriginDuplicateFunc func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error)
	CreateOriginFunc       func(ctx context.Context, origin *models.Origin) error
	UpdateOriginFunc       func(ctx context.Context, origin *models.Origin) error
	DeleteOriginFunc       func(ctx context.Context, origin *models.Origin) error
}

func (r *mockOriginRepository) GetOrigins(ctx context.Context) ([]models.Origin, error) {
	if r.GetOriginsFunc != nil {
		return r.GetOriginsFunc(ctx)
	}

	return nil, nil
}

func (r *mockOriginRepository) GetActiveOrigins(ctx context.Context) ([]models.Origin, error) {
	if r.GetActiveOriginsFunc != nil {
		return r.GetActiveOriginsFunc(ctx)
	}
	return nil, nil
}

func (r *mockOriginRepository) GetOriginByID(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
	if r.GetOriginByIDFunc != nil {
		return r.GetOriginByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (r *mockOriginRepository) GetOriginsByName(ctx context.Context, name, lang string) ([]models.Origin, error) {
	if r.GetOriginsByNameFunc != nil {
		return r.GetOriginsByNameFunc(ctx, name, lang)
	}
	return nil, nil
}

func (r *mockOriginRepository) GetOriginDuplicate(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
	if r.GetOriginDuplicateFunc != nil {
		return r.GetOriginDuplicateFunc(ctx, names, ID)
	}
	return nil, nil
}

func (r *mockOriginRepository) CreateOrigin(ctx context.Context, origin *models.Origin) error {
	if r.CreateOriginFunc != nil {
		return r.CreateOriginFunc(ctx, origin)
	}
	return nil
}

func (r *mockOriginRepository) UpdateOrigin(ctx context.Context, origin *models.Origin) error {
	if r.UpdateOriginFunc != nil {
		return r.UpdateOriginFunc(ctx, origin)
	}
	return nil
}

func (r *mockOriginRepository) DeleteOrigin(ctx context.Context, origin *models.Origin) error {
	if r.DeleteOriginFunc != nil {
		return r.DeleteOriginFunc(ctx, origin)
	}
	return nil
}

func TestOriginFindAll(t *testing.T) {
	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)

	t.Run("Success", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
				return []models.Origin{
					origin,
				}, nil
			},
		}

		service := services.NewOriginService(&originRepo)
		expected := []models.Origin{origin}

		origins, err := service.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, origins)
	})

	t.Run("Error", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewOriginService(&originRepo)
		origins, err := service.FindAll(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, origins)
	})
}

func TestOriginFindAllActive(t *testing.T) {
	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)

	t.Run("Success", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetActiveOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
				return []models.Origin{
					origin,
				}, nil
			},
		}

		service := services.NewOriginService(&originRepo)
		expected := []models.OriginFormResponse{origin.ToFormResponse("pt")}

		origins, err := service.FindAllActive(context.Background(), "pt")

		assert.NoError(t, err)
		assert.Equal(t, expected, origins)
	})

	t.Run("Error", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetActiveOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewOriginService(&originRepo)
		origins, err := service.FindAllActive(context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, origins)
	})
}

func TestOriginFindByID(t *testing.T) {
	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)

	t.Run("Success", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &origin, nil
			},
		}

		service := services.NewOriginService(&originRepo)
		originFound, err := service.FindByID(context.Background(), origin.ID)

		assert.NoError(t, err)
		assert.Equal(t, &origin, originFound)
	})

	t.Run("Record not found", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewOriginService(&originRepo)
		origin, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, origin)
	})

	t.Run("DB error", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewOriginService(&originRepo)
		origin, err := service.FindByID(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, origin)
	})
}

func TestOriginFindByName(t *testing.T) {
	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)

	t.Run("Success", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginsByNameFunc: func(ctx context.Context, name, lang string) ([]models.Origin, error) {
				return []models.Origin{origin}, nil
			},
		}

		service := services.NewOriginService(&originRepo)
		origins, err := service.FindByName(context.Background(), "huma", "en")

		assert.NoError(t, err)
		assert.Equal(t, []models.Origin{origin}, origins)
	})

	t.Run("Error", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginsByNameFunc: func(ctx context.Context, name, lang string) ([]models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewOriginService(&originRepo)
		origins, err := service.FindByName(context.Background(), "huma", "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, origins)
	})
}

func TestOriginCreate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		originRepo := mockOriginRepository{
			CreateOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return nil
			},
		}

		service := services.NewOriginService(&originRepo)
		err := service.Create(context.Background(), &models.Origin{})

		assert.NoError(t, err)
	})

	t.Run("Error - Find duplicate", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewOriginService(&originRepo)
		err := service.Create(context.Background(), &models.Origin{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
		}

		service := services.NewOriginService(&originRepo)
		err := service.Create(context.Background(), &models.Origin{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
	})

	t.Run("Error - Create", func(t *testing.T) {
		originRepo := mockOriginRepository{
			CreateOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewOriginService(&originRepo)
		err := service.Create(context.Background(), &models.Origin{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
	})
}

func TestOriginUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {})

	t.Run("Error - Not Found", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewOriginService(&originRepo)
		origin, err := service.Update(context.Background(), uuid.New(), models.OriginUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, origin)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
			GetOriginDuplicateFunc: func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
		}

		service := services.NewOriginService(&originRepo)
		origin, err := service.Update(context.Background(), uuid.New(), models.OriginUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Empty(t, origin)
	})

	t.Run("Error - Update", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
			UpdateOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewOriginService(&originRepo)
		origin, err := service.Update(context.Background(), uuid.New(), models.OriginUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, origin)
	})
}

func TestOriginDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
			DeleteOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return nil
			},
		}

		service := services.NewOriginService(&originRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewOriginService(&originRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("Error - Delete", func(t *testing.T) {
		originRepo := mockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
				return &models.Origin{}, nil
			},
			DeleteOriginFunc: func(ctx context.Context, origin *models.Origin) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewOriginService(&originRepo)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
	})
}
