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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestAdminUserFind(t *testing.T) {
	user := testmodels.NewAdminLoginUser()
	lang := "en"

	userResponse := user.ToAdminResponse(lang)

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context, filter models.AdminUserFilter) ([]models.User, error) {
				return []models.User{user}, nil
			},
		}

		service := services.NewAdminUserService(userRepo, nil, nil, nil)
		result, err := service.Find(
			context.Background(), models.AdminUserFilter{}, lang)

		expected := []models.AdminUserResponse{userResponse}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context, filter models.AdminUserFilter) ([]models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.Find(context.Background(), models.AdminUserFilter{}, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminUserFindByID(t *testing.T) {
	user := testmodels.NewAdminLoginUser()
	lang := "en"

	userResponse := user.ToAdminResponse(lang)

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
		}

		service := services.NewAdminUserService(userRepo, nil, nil, nil)
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

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
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

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.FindByID(context.Background(), uuid.New(), lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminUserFindByUsername(t *testing.T) {
	user := testmodels.NewAdminLoginUser()
	lang := "en"

	userResponse := user.ToAdminResponse(lang)

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return &user, nil
			},
		}

		service := services.NewAdminUserService(userRepo, nil, nil, nil)
		result, err := service.FindByUsername(context.Background(), user.Username, lang)

		assert.NoError(t, err)
		assert.Equal(t, &userResponse, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.FindByUsername(context.Background(), "invalid", lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.FindByUsername(context.Background(), "invalid", lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminUserFindByEmail(t *testing.T) {
	user := testmodels.NewAdminLoginUser()
	lang := "en"

	userResponse := user.ToAdminResponse(lang)

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
				return &user, nil
			},
		}

		service := services.NewAdminUserService(userRepo, nil, nil, nil)
		result, err := service.FindByEmail(context.Background(), user.Email, lang)

		assert.NoError(t, err)
		assert.Equal(t, &userResponse, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.FindByEmail(context.Background(), "invalid@mail.com", lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.FindByEmail(context.Background(), "invalid@mail.com", lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminUserCreate(t *testing.T) {
	lang := "en"
	adminName := "Roberto"

	input := models.AdminUserCreateInput{
		Name:        "Admin",
		Username:    "admin",
		Email:       "admin@mail.com",
		Password:    "123456",
		CountryCode: "BRA",
	}

	country := testmodels.NewCountry("BRA", map[string]string{"en": "Brazil"})

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				return nil
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &country, nil
			},
		}

		hasher := &mocks.MockHasher{}

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, nil)
		result, err := service.Create(context.Background(), input, adminName, lang)
		result.ActivatedOn = nil

		expected := models.AdminUserResponse{
			Name:        input.Name,
			Username:    input.Username,
			Email:       input.Email,
			CountryCode: input.CountryCode,
			Country:     country.Names[lang],
			ActivatedBy: &adminName,
			CreatedBy:   adminName,
		}

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Error - Conflict Username", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return &models.User{}, nil
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.Create(context.Background(), input, adminName, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflictUsername)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict Email", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return &models.User{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.Create(context.Background(), input, adminName, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflictEmail)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
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

		hasher := &mocks.MockHasher{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, mockLogger)
		result, err := service.Create(context.Background(), input, adminName, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCountryCode)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
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

		hasher := &mocks.MockHasher{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, mockLogger)
		result, err := service.Create(context.Background(), input, adminName, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Password hash failed", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &country, nil
			},
		}

		hasher := &mocks.MockHasher{
			HashFunc: func(password string) (string, error) {
				return "", bcrypt.ErrPasswordTooLong
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		service := services.NewAdminUserService(userRepo, countryRepo, hasher, mockLogger)
		result, err := service.Create(context.Background(), input, adminName, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &country, nil
			},
		}

		hasher := &mocks.MockHasher{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, mockLogger)
		result, err := service.Create(context.Background(), input, adminName, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminUserUpdate(t *testing.T) {
	lang := "en"
	userID := uuid.New()

	existingUser := testmodels.NewAdminLoginUser()
	existingUser.ID = userID

	newUsername := "admin2"
	newEmail := "admin2@mail.com"
	newPassword := "newpassword"
	countryCode := "BRA"

	input := models.AdminUserUpdateInput{
		Username:    &newUsername,
		Email:       &newEmail,
		Password:    &newPassword,
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

		hasher := &mocks.MockHasher{}

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, nil)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
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

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
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

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflictUsername)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict Email", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return &models.User{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflictEmail)
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

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Duplicate Email Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
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

		hasher := &mocks.MockHasher{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, mockLogger)
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

		hasher := &mocks.MockHasher{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Password hash failed", func(t *testing.T) {
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
				return &country, nil
			},
		}

		hasher := &mocks.MockHasher{
			HashFunc: func(password string) (string, error) {
				return "", bcrypt.ErrPasswordTooLong
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, mockLogger)
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

		hasher := &mocks.MockHasher{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, countryRepo, hasher, mockLogger)
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminActivateUser(t *testing.T) {
	user := testmodels.NewLoginUser()
	adminName := "admin"

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
		}

		user.IsActive = false
		service := services.NewAdminUserService(userRepo, nil, nil, nil)
		err := service.ActivateUser(context.Background(), user.ID, adminName)

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		err := service.ActivateUser(context.Background(), user.ID, adminName)

		assert.Error(t, err)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			UpdateUserFunc: func(ctx context.Context, user *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		user.IsActive = false
		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		err := service.ActivateUser(context.Background(), user.ID, adminName)

		assert.Error(t, err)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminDeactivateUser(t *testing.T) {
	user := testmodels.NewLoginUser()

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
		}

		user.IsActive = true
		service := services.NewAdminUserService(userRepo, nil, nil, nil)
		err := service.DeactivateUser(context.Background(), user.ID)

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		err := service.DeactivateUser(context.Background(), user.ID)

		assert.Error(t, err)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			UpdateUserFunc: func(ctx context.Context, user *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		user.IsActive = true
		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		err := service.DeactivateUser(context.Background(), user.ID)

		assert.Error(t, err)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminUserDelete(t *testing.T) {
	user := testmodels.NewAdminLoginUser()

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			DeleteUserFunc: func(ctx context.Context, user *models.User) error {
				return nil
			},
		}

		service := services.NewAdminUserService(userRepo, nil, nil, nil)
		err := service.Delete(context.Background(), user.ID)

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			DeleteUserFunc: func(ctx context.Context, user *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAdminUserService(userRepo, nil, nil, mockLogger)
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
