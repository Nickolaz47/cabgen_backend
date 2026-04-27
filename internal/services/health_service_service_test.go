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

func TestHealthServiceFindAll(t *testing.T) {
	country := testmodels.NewCountry("BRA", nil)
	healthService := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John Doe", "john@example.com", "123456789",
		true,
	)

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServicesFunc: func(ctx context.Context) (
				[]models.HealthService, error) {
				return []models.HealthService{healthService}, nil
			},
		}
		service := services.NewHealthServiceService(repo, nil, nil)

		expected := []models.HealthServiceAdminTableResponse{
			healthService.ToAdminTableResponse(),
		}
		healthServices, err := service.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, healthServices)
	})

	t.Run("Error", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServicesFunc: func(ctx context.Context) (
				[]models.HealthService, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		healthServices, err := service.FindAll(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, healthServices)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestHealthServiceFindAllActive(t *testing.T) {
	country := testmodels.NewCountry("BRA", nil)
	healthService := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John Doe", "john@example.com", "123456789",
		true,
	)

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetActiveHealthServicesFunc: func(ctx context.Context) (
				[]models.HealthService, error) {
				return []models.HealthService{healthService}, nil
			},
		}
		service := services.NewHealthServiceService(repo, nil, nil)

		expected := []models.HealthServiceFormResponse{
			healthService.ToFormResponse(),
		}
		healthServices, err := service.FindAllActive(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, healthServices)
	})

	t.Run("Error", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetActiveHealthServicesFunc: func(ctx context.Context) ([]models.HealthService, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		healthServices, err := service.FindAllActive(context.Background())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, healthServices)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestHealthServiceFindByID(t *testing.T) {
	country := testmodels.NewCountry("BRA", nil)
	healthService := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John Doe", "john@example.com", "123456789",
		true,
	)

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.HealthService, error) {
				return &healthService, nil
			},
		}
		service := services.NewHealthServiceService(repo, nil, nil)

		expected := healthService.ToAdminTableResponse()
		result, err := service.FindByID(context.Background(), healthService.ID)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Record Not Found", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		result, err := service.FindByID(context.Background(), healthService.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.HealthService, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		result, err := service.FindByID(context.Background(), healthService.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestHealthServiceFindByName(t *testing.T) {
	country := testmodels.NewCountry("BRA", nil)
	healthService := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John Doe", "john@example.com", "123456789",
		true,
	)

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServicesByNameFunc: func(ctx context.Context,
				name string) ([]models.HealthService, error) {
				return []models.HealthService{healthService}, nil
			},
		}
		service := services.NewHealthServiceService(repo, nil, nil)

		expected := []models.HealthServiceAdminTableResponse{
			healthService.ToAdminTableResponse(),
		}
		result, err := service.FindByName(context.Background(), "Hosp")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServicesByNameFunc: func(ctx context.Context, name string) (
				[]models.HealthService, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		result, err := service.FindByName(context.Background(), "Hosp")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestHealthServiceCreate(t *testing.T) {
	country := testmodels.NewCountry("", nil)
	input := models.HealthServiceCreateInput{
		Name:         "Hospital A",
		Type:         models.Public,
		CountryCode:  "BRA",
		City:         "Rio de Janeiro",
		Contactant:   "John Doe",
		ContactEmail: "john@example.com",
		ContactPhone: "123456789",
		IsActive:     true,
	}

	t.Run("Success", func(t *testing.T) {
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (
				*models.Country, error) {
				return &country, nil
			},
		}
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateHealthServiceFunc: func(ctx context.Context, healthService *models.HealthService) error {
				return nil
			},
		}
		service := services.NewHealthServiceService(repo, countryRepo, nil)

		expected := models.HealthServiceAdminTableResponse{
			Name:         input.Name,
			Type:         input.Type,
			Country:      country.Code,
			City:         input.City,
			Contactant:   input.Contactant,
			ContactEmail: input.ContactEmail,
			ContactPhone: input.ContactPhone,
			IsActive:     input.IsActive,
		}
		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Invalid Country Code", func(t *testing.T) {
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (
				*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		repo := &mocks.MockHealthServiceRepository{}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		service := services.NewHealthServiceService(repo, countryRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCountryCode)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Repo Internal", func(t *testing.T) {
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (
				*models.Country, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		repo := &mocks.MockHealthServiceRepository{}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		service := services.NewHealthServiceService(repo, countryRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{ID: 1, Code: "BRA"}, nil
			},
		}
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
				return &models.HealthService{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		service := services.NewHealthServiceService(repo, countryRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Create", func(t *testing.T) {
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (
				*models.Country, error) {
				return &country, nil
			},
		}

		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceDuplicateFunc: func(
				ctx context.Context, name string, ID uuid.UUID) (
					*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateHealthServiceFunc: func(
				ctx context.Context, healthService *models.HealthService) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		service := services.NewHealthServiceService(repo, countryRepo, mockLogger)
		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestHealthServiceUpdate(t *testing.T) {
	id := uuid.New()
	name := "Hospital B"
	countryCode := "USA"
	input := models.HealthServiceUpdateInput{
		Name:        &name,
		CountryCode: &countryCode,
	}

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
				return &models.HealthService{ID: id}, nil
			},
			GetHealthServiceDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
			UpdateHealthServiceFunc: func(ctx context.Context, healthService *models.HealthService) error {
				return nil
			},
		}
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{ID: 2, Code: "USA"}, nil
			},
		}

		service := services.NewHealthServiceService(repo, countryRepo, nil)
		result, err := service.Update(context.Background(), id, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		result, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
				return &models.HealthService{ID: id}, nil
			},
			GetHealthServiceDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
				return &models.HealthService{ID: uuid.New()}, nil
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		result, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflict)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Invalid Country Code", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
				return &models.HealthService{ID: id}, nil
			},
			GetHealthServiceDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, countryRepo, mockLogger)
		result, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCountryCode)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
				return &models.HealthService{ID: id}, nil
			},
			GetHealthServiceDuplicateFunc: func(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
			UpdateHealthServiceFunc: func(ctx context.Context, healthService *models.HealthService) error {
				return gorm.ErrInvalidTransaction
			},
		}
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{ID: 2, Code: "USA"}, nil
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, countryRepo, mockLogger)
		result, err := service.Update(context.Background(), id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestHealthServiceDelete(t *testing.T) {
	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
				return &models.HealthService{ID: id}, nil
			},
			DeleteHealthServiceFunc: func(ctx context.Context, healthService *models.HealthService) error {
				return nil
			},
		}

		service := services.NewHealthServiceService(repo, nil, nil)
		err := service.Delete(context.Background(), id)

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		err := service.Delete(context.Background(), id)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		repo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
				return &models.HealthService{ID: id}, nil
			},
			DeleteHealthServiceFunc: func(ctx context.Context, healthService *models.HealthService) error {
				return gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewHealthServiceService(repo, nil, mockLogger)
		err := service.Delete(context.Background(), id)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
