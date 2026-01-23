package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCountryFindAll(t *testing.T) {
	country := testmodels.NewCountry(
		"BRA",
		map[string]string{"pt": "Brasil", "en": "Brazil"},
	)

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountriesFunc: func(ctx context.Context) ([]models.Country, error) {
				return []models.Country{country}, nil
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.FindAll(context.Background(), "pt")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, country.ToFormResponse("pt"), result[0])
	})

	t.Run("Error", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountriesFunc: func(ctx context.Context) ([]models.Country, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.FindAll(context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
	})
}

func TestCountryFindByCode(t *testing.T) {
	country := testmodels.NewCountry(
		"BRA",
		map[string]string{"pt": "Brasil"},
	)

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &country, nil
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.FindByCode(context.Background(), "BRA")

		assert.NoError(t, err)
		assert.Equal(t, country.ToAdminDetailResponse(), *result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.FindByCode(context.Background(), "BRA")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
	})

	t.Run("Error - Internal", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.FindByCode(context.Background(), "BRA")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
	})
}

func TestCountriesFindByName(t *testing.T) {
	country := testmodels.NewCountry(
		"BRA",
		map[string]string{"en": "Brazil"},
	)

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountriesByNameFunc: func(ctx context.Context, name, lang string) ([]models.Country, error) {
				return []models.Country{country}, nil
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.FindByName(context.Background(), "bra", "en")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, country.ToFormResponse("en"), result[0])
	})

	t.Run("Error", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountriesByNameFunc: func(ctx context.Context, name, lang string) ([]models.Country, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.FindByName(context.Background(), "bra", "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
	})
}

func TestCountryCreate(t *testing.T) {
	input := models.CountryCreateInput{
		Code:  "BRA",
		Names: map[string]string{"pt": "Brasil"},
	}

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryDuplicateFunc: func(ctx context.Context, names models.JSONMap, code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateCountryFunc: func(ctx context.Context, country *models.Country) error {
				return nil
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryDuplicateFunc: func(ctx context.Context, names models.JSONMap, code string) (*models.Country, error) {
				return &models.Country{}, nil
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Nil(t, result)
	})

	t.Run("Error - Create", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryDuplicateFunc: func(ctx context.Context, names models.JSONMap, code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateCountryFunc: func(ctx context.Context, country *models.Country) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
	})
}

func TestCountryUpdate(t *testing.T) {
	newCode := "BRZ"
	input := models.CountryUpdateInput{
		Code:  &newCode,
		Names: map[string]string{"pt": "Brasil"},
	}

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{Code: "BRA"}, nil
			},
			UpdateCountryFunc: func(ctx context.Context, country *models.Country) error {
				return nil
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.Update(context.Background(), "BRA", models.CountryUpdateInput{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.Update(context.Background(), "BRA", models.CountryUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
	})

	t.Run("Error - Conflict (Code)", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{Code: "BRA"}, nil
			},
			GetCountryDuplicateFunc: func(ctx context.Context, names models.JSONMap, code string) (*models.Country, error) {
				return &models.Country{}, nil
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.Update(context.Background(), "BRA", input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Nil(t, result)
	})

	t.Run("Error - Update", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{Code: "BRA"}, nil
			},
			UpdateCountryFunc: func(ctx context.Context, country *models.Country) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewCountryService(repo)
		result, err := service.Update(context.Background(), "BRA", models.CountryUpdateInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
	})
}

func TestCountryDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{}, nil
			},
			DeleteCountryFunc: func(ctx context.Context, country *models.Country) error {
				return nil
			},
		}

		service := services.NewCountryService(repo)
		err := service.Delete(context.Background(), "BRA")

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewCountryService(repo)
		err := service.Delete(context.Background(), "BRA")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("Error - Delete", func(t *testing.T) {
		repo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{}, nil
			},
			DeleteCountryFunc: func(ctx context.Context, country *models.Country) error {
				return gorm.ErrInvalidTransaction
			},
		}

		service := services.NewCountryService(repo)
		err := service.Delete(context.Background(), "BRA")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
	})
}
