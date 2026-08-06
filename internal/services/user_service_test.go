package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
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
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
		}

		service := services.NewUserService(userRepo, nil, nil, nil, nil, nil, "")
		result, err := service.FindByID(context.Background(), user.ID, lang)

		assert.NoError(t, err)
		assert.Equal(t, &userResponse, result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		result, err := service.FindByID(context.Background(), uuid.New(), lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		result, err := service.FindByID(context.Background(), uuid.New(), lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestUserDelete(t *testing.T) {
	user := testmodels.NewLoginUser()

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			DeleteUserFunc: func(ctx context.Context,
				user *models.User) error {
				return nil
			},
		}

		service := services.NewUserService(userRepo, nil, nil, nil,
			&mocks.MockTaskEnqueuer{}, nil, t.TempDir())
		err := service.Delete(context.Background(), user.ID)

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			DeleteUserFunc: func(ctx context.Context, user *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		err := service.Delete(context.Background(), uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
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
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			UpdateUserFunc: func(ctx context.Context, user *models.User) error {
				return nil
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context,
				code string) (*models.Country, error) {
				return &country, nil
			},
		}

		service := services.NewUserService(userRepo, countryRepo, nil, nil, nil,
			nil, "")
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Get User Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Conflict Username", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string,
				ID uuid.UUID) (*models.User, error) {
				return &models.User{}, nil
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflictUsername)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Duplicate Username Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context,
				code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, countryRepo, nil, nil, nil,
			mockLogger, "")
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCountryCode)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context,
				code string) (*models.Country, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, countryRepo, nil, nil, nil,
			mockLogger, "")
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update Internal Server", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &existingUser, nil
			},
			ExistsByUsernameFunc: func(ctx context.Context, username *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			ExistsByEmailFunc: func(ctx context.Context, email *string,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			UpdateUserFunc: func(ctx context.Context, user *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}

		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context,
				code string) (*models.Country, error) {
				return &country, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, countryRepo, nil, nil, nil,
			mockLogger, "")
		result, err := service.Update(context.Background(), userID, input, lang)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestUpdatePassword(t *testing.T) {
	user := testmodels.NewLoginUser()
	user.Password = "$2a$10$fakehashedpassword"

	input := models.UpdatePasswordInput{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			UpdateUserFunc: func(ctx context.Context, u *models.User) error {
				return nil
			},
		}
		hasher := &mocks.MockPasswordHasher{
			CheckPasswordFunc: func(hashPassword, password string) error {
				return nil
			},
			HashFunc: func(password string) (string, error) {
				return "$2a$10$newhashedpassword", nil
			},
		}

		service := services.NewUserService(userRepo, nil, nil, hasher, nil,
			nil, "")
		err := service.UpdatePassword(context.Background(), user.ID, input)

		assert.NoError(t, err)
	})

	t.Run("Error - Wrong Current Password", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
		}
		hasher := &mocks.MockPasswordHasher{
			CheckPasswordFunc: func(hashPassword, password string) error {
				return errors.New(
					"crypto/bcrypt: hashedPassword is not the hash of the" +
						" given password")
			},
		}
		mockLogger, _ := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, hasher, nil,
			mockLogger, "")
		err := service.UpdatePassword(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrCurrentPasswordMismatch)
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, _ := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		err := service.UpdatePassword(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
	})
}

func TestRequestEmailUpdate(t *testing.T) {
	user := testmodels.NewLoginUser()
	newEmail := "new@example.com"

	input := models.RequestEmailUpdateInput{
		NewEmail:        newEmail,
		ConfirmNewEmail: newEmail,
	}

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			GetUserByEmailFunc: func(ctx context.Context,
				email string) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		var deletedUserID uuid.UUID
		var createdReq *models.EmailUpdateRequest
		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			DeleteRequestsByUserIDFunc: func(ctx context.Context,
				userID uuid.UUID) error {
				deletedUserID = userID
				return nil
			},
			CreateRequestFunc: func(ctx context.Context,
				req *models.EmailUpdateRequest) error {
				createdReq = req
				return nil
			},
		}

		var enqueuedTaskType string
		taskEnqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task,
				opts ...asynq.Option) (*asynq.TaskInfo, error) {
				enqueuedTaskType = task.Type()
				return &asynq.TaskInfo{ID: "task-123"}, nil
			},
		}

		service := services.NewUserService(userRepo, nil, emailUpdateRepo, nil,
			taskEnqueuer, nil, "")
		err := service.RequestEmailUpdate(context.Background(), user.ID, input)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, deletedUserID)
		assert.NotNil(t, createdReq)
		assert.Equal(t, user.ID, createdReq.UserID)
		assert.Equal(t, user.Email, createdReq.OldEmail)
		assert.Equal(t, newEmail, createdReq.NewEmail)
		assert.NotEmpty(t, createdReq.Token)
		assert.WithinDuration(t, time.Now().Add(15*time.Minute),
			createdReq.ExpiresAt, time.Second)
		assert.Equal(t, tasks.TaskTypeEmailUpdateConfirmation, enqueuedTaskType)
	})

	t.Run("Error - Create Request Internal", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			GetUserByEmailFunc: func(ctx context.Context,
				email string) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			CreateRequestFunc: func(ctx context.Context,
				req *models.EmailUpdateRequest) error {
				return gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, emailUpdateRepo, nil,
			nil, mockLogger, "")
		err := service.RequestEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Enqueue Internal", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			GetUserByEmailFunc: func(ctx context.Context,
				email string) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			CreateRequestFunc: func(ctx context.Context,
				req *models.EmailUpdateRequest) error {
				return nil
			},
		}
		taskEnqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task,
				opts ...asynq.Option) (*asynq.TaskInfo, error) {
				return nil, errors.New("redis error")
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, emailUpdateRepo, nil,
			taskEnqueuer, mockLogger, "")
		err := service.RequestEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Duplicate Task", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			GetUserByEmailFunc: func(ctx context.Context,
				email string) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			CreateRequestFunc: func(ctx context.Context,
				req *models.EmailUpdateRequest) error {
				return nil
			},
		}
		taskEnqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task,
				opts ...asynq.Option) (*asynq.TaskInfo, error) {
				return nil, asynq.ErrDuplicateTask
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, emailUpdateRepo, nil,
			taskEnqueuer, mockLogger, "")
		err := service.RequestEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrDuplicateTask)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, _ := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		err := service.RequestEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("Error - Email Same", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
		}

		sameEmailInput := models.RequestEmailUpdateInput{
			NewEmail:        user.Email,
			ConfirmNewEmail: user.Email,
		}

		service := services.NewUserService(userRepo, nil, nil, nil, nil, nil, "")
		err := service.RequestEmailUpdate(context.Background(), user.ID,
			sameEmailInput)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrEmailSame)
	})

	t.Run("Error - Email Already Exists", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			GetUserByEmailFunc: func(ctx context.Context,
				email string) (*models.User, error) {
				return &models.User{ID: uuid.New()}, nil
			},
		}

		service := services.NewUserService(userRepo, nil, nil, nil, nil, nil, "")
		err := service.RequestEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrConflictEmail)
	})

	t.Run("Error - Get User By Email Internal", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			GetUserByEmailFunc: func(ctx context.Context,
				email string) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, nil, nil, nil,
			mockLogger, "")
		err := service.RequestEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestConfirmEmailUpdate(t *testing.T) {
	user := testmodels.NewLoginUser()
	reqID := uuid.New()
	token := "valid-token"
	newEmail := "new@example.com"

	input := models.ConfirmEmailUpdateInput{Token: token}
	request := models.EmailUpdateRequest{
		ID:        reqID,
		UserID:    user.ID,
		OldEmail:  user.Email,
		NewEmail:  newEmail,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			UpdateUserFunc: func(ctx context.Context, u *models.User) error {
				return nil
			},
		}
		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			GetByTokenFunc: func(ctx context.Context,
				t string) (*models.EmailUpdateRequest, error) {
				return &request, nil
			},
		}

		service := services.NewUserService(userRepo, nil, emailUpdateRepo, nil,
			nil, nil, "")
		err := service.ConfirmEmailUpdate(context.Background(), user.ID, input)

		assert.NoError(t, err)
	})

	t.Run("Error - Invalid Token", func(t *testing.T) {
		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			GetByTokenFunc: func(ctx context.Context,
				t string) (*models.EmailUpdateRequest, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		service := services.NewUserService(nil, nil, emailUpdateRepo, nil, nil,
			nil, "")
		err := service.ConfirmEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidEmailUpdateToken)
	})

	t.Run("Error - Token User Mismatch", func(t *testing.T) {
		otherRequest := request
		otherRequest.UserID = uuid.New()

		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			GetByTokenFunc: func(ctx context.Context,
				t string) (*models.EmailUpdateRequest, error) {
				return &otherRequest, nil
			},
		}

		service := services.NewUserService(nil, nil, emailUpdateRepo, nil, nil,
			nil, "")
		err := service.ConfirmEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidEmailUpdateToken)
	})

	t.Run("Error - Expired Token", func(t *testing.T) {
		expiredRequest := request
		expiredRequest.ExpiresAt = time.Now().Add(-15 * time.Minute)

		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			GetByTokenFunc: func(ctx context.Context,
				t string) (*models.EmailUpdateRequest, error) {
				return &expiredRequest, nil
			},
		}

		service := services.NewUserService(nil, nil, emailUpdateRepo, nil, nil,
			nil, "")
		err := service.ConfirmEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrExpiredEmailUpdateToken)
	})

	t.Run("Error - Get User Internal", func(t *testing.T) {
		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			GetByTokenFunc: func(ctx context.Context,
				t string) (*models.EmailUpdateRequest, error) {
				return &request, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, emailUpdateRepo, nil,
			nil, mockLogger, "")
		err := service.ConfirmEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update User Internal", func(t *testing.T) {
		emailUpdateRepo := &mocks.MockEmailUpdateRepository{
			GetByTokenFunc: func(ctx context.Context,
				t string) (*models.EmailUpdateRequest, error) {
				return &request, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &user, nil
			},
			UpdateUserFunc: func(ctx context.Context, u *models.User) error {
				return gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewUserService(userRepo, nil, emailUpdateRepo, nil,
			nil, mockLogger, "")
		err := service.ConfirmEmailUpdate(context.Background(), user.ID, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestEmailUpdateRequestIsExpired(t *testing.T) {
	t.Run("Expired", func(t *testing.T) {
		req := models.EmailUpdateRequest{
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		}

		assert.True(t, req.IsExpired())
	})

	t.Run("Not Expired", func(t *testing.T) {
		req := models.EmailUpdateRequest{
			ExpiresAt: time.Now().Add(1 * time.Minute),
		}

		assert.False(t, req.IsExpired())
	})
}
