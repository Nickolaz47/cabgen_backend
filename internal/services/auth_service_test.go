package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	input := testmodels.NewRegisterUser("", "")
	lang := "en"

	validCountry := testmodels.NewCountry("", nil)

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				return nil
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &validCountry, nil
			},
		}

		hasher := &mocks.MockHasher{}
		enqueuer := &mocks.MockTaskEnqueuer{}
		mockLogger, logs := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewAuthService(
			userRepo, countryRepo, hasher, nil, enqueuer, mockLogger)

		expected := models.UserResponse{
			Name:        input.Name,
			Username:    input.Username,
			Email:       input.Email,
			CountryCode: input.CountryCode,
			Country:     validCountry.Names[lang],
			UserRole:    models.Collaborator,
			Interest:    input.Interest,
			Role:        input.Role,
			Institution: input.Institution,
		}
		result, err := svc.Register(ctx, input, lang)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Email already exists", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByEmailFunc: func(ctx context.Context, email *string, _ uuid.UUID) (*models.User, error) {
				return &models.User{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAuthService(userRepo, nil, nil, nil, nil, mockLogger)
		result, err := svc.Register(ctx, input, lang)

		assert.Equal(t, services.ErrConflictEmail, err)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - ExistsByEmail Internal", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByEmailFunc: func(ctx context.Context, email *string, _ uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAuthService(userRepo, nil, nil, nil, nil, mockLogger)
		result, err := svc.Register(ctx, input, lang)

		assert.Equal(t, services.ErrInternal, err)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Username already exists", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByEmailFunc: func(ctx context.Context, email *string, _ uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, _ uuid.UUID) (*models.User, error) {
				return &models.User{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAuthService(userRepo, nil, nil, nil, nil, mockLogger)
		result, err := svc.Register(ctx, input, lang)

		assert.Equal(t, services.ErrConflictUsername, err)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - ExistsByUsername Internal", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByEmailFunc: func(ctx context.Context, email *string, _ uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, _ uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAuthService(userRepo, nil, nil, nil, nil, mockLogger)
		result, err := svc.Register(ctx, input, lang)

		assert.Equal(t, services.ErrInternal, err)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Email Mismatch", func(t *testing.T) {
		badInput := input
		badInput.ConfirmEmail = "wrong@mail.com"

		userRepo := &mocks.MockUserRepository{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAuthService(userRepo, nil, nil, nil, nil, mockLogger)
		result, err := svc.Register(ctx, badInput, lang)

		assert.Equal(t, services.ErrEmailMismatch, err)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Password Mismatch", func(t *testing.T) {
		badInput := input
		badInput.ConfirmPassword = "wrong"

		userRepo := &mocks.MockUserRepository{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAuthService(userRepo, nil, nil, nil, nil, mockLogger)
		result, err := svc.Register(ctx, badInput, lang)

		assert.Equal(t, services.ErrPasswordMismatch, err)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Hash Internal", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return &models.User{}, nil
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &models.Country{}, nil
			},
		}

		hasher := &mocks.MockHasher{
			HashFunc: func(password string) (string, error) {
				return "", bcrypt.ErrHashTooShort
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAuthService(userRepo, countryRepo, hasher, nil, nil, mockLogger)
		result, err := svc.Register(ctx, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Invalid Country Code", func(t *testing.T) {
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

		svc := services.NewAuthService(userRepo, countryRepo, hasher, nil, nil, mockLogger)
		result, err := svc.Register(ctx, input, lang)

		assert.Equal(t, services.ErrInvalidCountryCode, err)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - GetCountryByCode Internal", func(t *testing.T) {
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

		svc := services.NewAuthService(userRepo, countryRepo, hasher, nil, nil, mockLogger)
		result, err := svc.Register(ctx, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &validCountry, nil
			},
		}

		hasher := &mocks.MockHasher{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, countryRepo, hasher, nil, nil, mockLogger)
		_, err := svc.Register(ctx, input, lang)

		assert.Equal(t, services.ErrInternal, err)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Success - Soft Fail Asynq", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			ExistsByEmailFunc: func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				return nil
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context, code string) (*models.Country, error) {
				return &validCountry, nil
			},
		}

		hasher := &mocks.MockHasher{}

		failingEnqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
				return nil, errors.New("redis timeout")
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, countryRepo, hasher, nil, failingEnqueuer, mockLogger)

		expected := models.UserResponse{
			Name:        input.Name,
			Username:    input.Username,
			Email:       input.Email,
			CountryCode: input.CountryCode,
			Country:     validCountry.Names[lang],
			UserRole:    models.Collaborator,
			Interest:    input.Interest,
			Role:        input.Role,
			Institution: input.Institution,
		}

		result, err := svc.Register(ctx, input, lang)

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	mockUser := testmodels.NewLoginUser()

	input := models.LoginInput{
		Username: mockUser.Username,
		Password: "12345678",
	}

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return &mockUser, nil
			},
		}

		hasher := &mocks.MockHasher{}
		provider := &mocks.MockTokenProvider{}

		svc := services.NewAuthService(userRepo, nil, hasher, provider, nil, nil)
		result, err := svc.Login(ctx, input)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Error - Username Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		hasher := &mocks.MockHasher{}
		provider := &mocks.MockTokenProvider{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, nil, hasher, provider, nil, mockLogger)
		result, err := svc.Login(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCredentials)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - GetUserByUsername Internal", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		hasher := &mocks.MockHasher{}
		provider := &mocks.MockTokenProvider{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, nil, hasher, provider, nil, mockLogger)
		result, err := svc.Login(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Disabled User", func(t *testing.T) {
		disabledUser := mockUser
		disabledUser.IsActive = false
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return &disabledUser, nil
			},
		}

		hasher := &mocks.MockHasher{}
		provider := &mocks.MockTokenProvider{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, nil, hasher, provider, nil, mockLogger)
		result, err := svc.Login(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrDisabledUser)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Wrong Password", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return &mockUser, nil
			},
		}

		hasher := &mocks.MockHasher{
			CheckPasswordFunc: func(hashPassword, password string) error {
				return bcrypt.ErrMismatchedHashAndPassword
			},
		}
		provider := &mocks.MockTokenProvider{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, nil, hasher, provider, nil, mockLogger)
		result, err := svc.Login(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCredentials)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - CheckPassword Internal", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return &mockUser, nil
			},
		}

		hasher := &mocks.MockHasher{
			CheckPasswordFunc: func(hashPassword, password string) error {
				return bcrypt.ErrPasswordTooLong
			},
		}
		provider := &mocks.MockTokenProvider{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, nil, hasher, provider, nil, mockLogger)
		result, err := svc.Login(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Generate Access Token", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return &mockUser, nil
			},
		}

		hasher := &mocks.MockHasher{}
		provider := &mocks.MockTokenProvider{
			GenerateTokenFunc: func(user models.UserToken, secret []byte, expiresIn time.Duration) (string, error) {
				return "", jwt.ErrInvalidKey
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, nil, hasher, provider, nil, mockLogger)
		result, err := svc.Login(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Generate Refresh Token", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
				return &mockUser, nil
			},
		}

		hasher := &mocks.MockHasher{}
		calls := 0
		provider := &mocks.MockTokenProvider{
			GenerateTokenFunc: func(user models.UserToken, secret []byte, expiresIn time.Duration) (string, error) {
				calls++
				if calls == 1 {
					return "", nil
				}
				return "", jwt.ErrInvalidKey
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(userRepo, nil, hasher, provider, nil, mockLogger)
		result, err := svc.Login(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestRefresh(t *testing.T) {
	ctx := context.Background()
	userToken := testmodels.NewUserToken(
		uuid.New(), "nikol47", models.Collaborator,
	)
	tokenStr := "requestToken"

	t.Run("Success", func(t *testing.T) {
		provider := &mocks.MockTokenProvider{
			ValidateTokenFunc: func(tokenStr string, secret []byte) (*models.UserToken, error) {
				return &userToken, nil
			},
			GenerateTokenFunc: func(user models.UserToken, secret []byte, expiresIn time.Duration) (string, error) {
				return "refreshedToken", nil
			},
		}

		svc := services.NewAuthService(nil, nil, nil, provider, nil, nil)
		result, err := svc.Refresh(ctx, tokenStr)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Error - ValidateToken Unauthorized", func(t *testing.T) {
		provider := &mocks.MockTokenProvider{
			ValidateTokenFunc: func(tokenStr string, secret []byte) (*models.UserToken, error) {
				return nil, errors.New("invalid token")
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(nil, nil, nil, provider, nil, mockLogger)
		result, err := svc.Refresh(ctx, tokenStr)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - GenerateToken Internal", func(t *testing.T) {
		provider := &mocks.MockTokenProvider{
			ValidateTokenFunc: func(tokenStr string, secret []byte) (*models.UserToken, error) {
				return &userToken, nil
			},
			GenerateTokenFunc: func(user models.UserToken, secret []byte, expiresIn time.Duration) (string, error) {
				return "", jwt.ErrInvalidKey
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAuthService(nil, nil, nil, provider, nil, mockLogger)
		result, err := svc.Refresh(ctx, tokenStr)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}
