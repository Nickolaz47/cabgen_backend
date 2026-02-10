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

func TestUserFindByID(t *testing.T) {
	user := testmodels.NewLoginUser()
	lang := "en"

	userResponse := user.ToResponse(lang)

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
		}

		service := services.NewUserService(userRepo, nil, nil)
		result, err := service.FindByID(context.Background(), user.ID, lang)

		assert.NoError(t, err)
		assert.Equal(t, &userResponse, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, mockLogger)
		result, err := service.FindByID(context.Background(), uuid.New(), lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, mockLogger)
		result, err := service.FindByID(context.Background(), uuid.New(), lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestUserUpdate(t *testing.T) {
	lang := "en"
	userID := uuid.New()

	existingUser := testmodels.NewLoginUser()
	existingUser.ID = userID

	newUsername := "john"
	countryCode := "BRA"

	input := models.UserUpdateInput{
		Username:    &newUsername,
		CountryCode: &countryCode,
	}

	country := testmodels.NewCountry("BRA", map[string]string{"en": "Brazil"})

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			UpdateUserFunc: func(ctx context.Context, user *models.User) error {
				return nil
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &country, nil
			},
		}

		service := services.NewUserService(userRepo, countryRepo, nil)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Get User Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict Username", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return &models.User{}, nil
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflictUsername)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Duplicate Username Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, countryRepo, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCountryCode)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, countryRepo, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			UpdateUserFunc: func(ctx context.Context, user *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &country, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, countryRepo, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}
