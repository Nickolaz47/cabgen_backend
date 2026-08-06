package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService interface {
	FindByID(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.UserUpdateInput, language string) (*models.UserResponse, error)
	Delete(ctx context.Context, ID uuid.UUID) error
	UpdatePassword(ctx context.Context, ID uuid.UUID, input models.UpdatePasswordInput) error
	RequestEmailUpdate(ctx context.Context, ID uuid.UUID, input models.RequestEmailUpdateInput) error
	ConfirmEmailUpdate(ctx context.Context, ID uuid.UUID, input models.ConfirmEmailUpdateInput) error
}

type userService struct {
	Repo            repositories.UserRepository
	CountryRepo     repositories.CountryRepository
	EmailUpdateRepo repositories.EmailUpdateRepository
	Hasher          security.PasswordHasher
	AsynqClient     TaskEnqueuer
	Logger          *zap.Logger
	RootDir         string
}

func NewUserService(
	repo repositories.UserRepository,
	countryRepo repositories.CountryRepository,
	emailUpdateRepo repositories.EmailUpdateRepository,
	hasher security.PasswordHasher,
	asynqClient TaskEnqueuer,
	logger *zap.Logger,
	rootDir string,
) UserService {
	return &userService{
		Repo:            repo,
		CountryRepo:     countryRepo,
		EmailUpdateRepo: emailUpdateRepo,
		Hasher:          hasher,
		AsynqClient:     asynqClient,
		Logger:          logger,
		RootDir:         rootDir,
	}
}

func (s *userService) FindByID(
	ctx context.Context, ID uuid.UUID,
	language string) (*models.UserResponse, error) {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "FindByID",
			logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "FindByID",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	response := user.ToResponse(language)
	return &response, nil
}

func (s *userService) Update(
	ctx context.Context, ID uuid.UUID,
	input models.UserUpdateInput,
	language string) (*models.UserResponse, error) {
	existingUser, err := s.Repo.GetUserByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "Update",
			logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "Update",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if input.Username != nil {
		duplicate, err := s.Repo.ExistsByUsername(
			ctx,
			input.Username,
			ID,
		)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"UserService", "Update",
				logging.DatabaseError, err,
			)...)
			return nil, ErrInternal
		}
		if duplicate != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"UserService", "Update",
				logging.DatabaseConflictUsernameError, err,
			)...)
			return nil, ErrConflictUsername
		}
	}

	if input.CountryCode != nil {
		country, err := s.CountryRepo.GetCountryByCode(ctx, *input.CountryCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error", logging.ServiceLogging(
					"UserService", "Update",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
				return nil, ErrInvalidCountryCode
			}
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"UserService", "Update",
				logging.ExternalRepositoryError, err,
			)...)
			return nil, ErrInternal
		}
		existingUser.CountryID = country.ID
		existingUser.Country = *country
	}

	validations.ApplyUpdateToUser(existingUser, &input)

	if err := s.Repo.UpdateUser(ctx, existingUser); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "Update",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	response := existingUser.ToResponse(language)
	return &response, nil
}

func (s *userService) Delete(ctx context.Context, ID uuid.UUID) error {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "Delete",
			logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "Delete",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	userEmail := user.Email
	userName := user.Name

	if err := s.Repo.DeleteUser(ctx, user); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "Delete",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	userFolder := filepath.Join(s.RootDir, "uploads", "users", user.ID.String())
	go func() {
		_ = os.RemoveAll(userFolder)
	}()

	if userEmail != "" {
		task, err := tasks.NewUserDeletedEmailTask(userEmail, userName)
		if err == nil {
			_, _ = s.AsynqClient.EnqueueContext(ctx, task,
				asynq.Queue(tasks.QueueEmail))
		}
	}

	return nil
}

func (s *userService) UpdatePassword(ctx context.Context, ID uuid.UUID,
	input models.UpdatePasswordInput) error {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"UserService", "UpdatePassword",
				logging.DatabaseNotFoundError, err,
			)...)
			return ErrNotFound
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "UpdatePassword", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if err := s.Hasher.CheckPassword(user.Password,
		input.CurrentPassword); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "UpdatePassword", logging.HasherError, err,
		)...)
		return ErrCurrentPasswordMismatch
	}

	hashedPassword, err := s.Hasher.Hash(input.NewPassword)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "UpdatePassword", logging.HasherError, err,
		)...)
		return ErrInternal
	}

	user.Password = hashedPassword
	if err := s.Repo.UpdateUser(ctx, user); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "UpdatePassword", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	return nil
}

func (s *userService) RequestEmailUpdate(ctx context.Context, ID uuid.UUID,
	input models.RequestEmailUpdateInput) error {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"UserService", "RequestEmailUpdate",
				logging.DatabaseNotFoundError, err,
			)...)
			return ErrNotFound
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "RequestEmailUpdate", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if input.NewEmail == user.Email {
		return ErrEmailSame
	}

	existing, err := s.Repo.GetUserByEmail(ctx, input.NewEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "RequestEmailUpdate", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}
	if existing != nil && existing.ID != ID {
		return ErrConflictEmail
	}

	if err := s.EmailUpdateRepo.DeleteRequestsByUserID(ctx, ID); err != nil {
		err = fmt.Errorf("%s: %v", ID, err)
		s.Logger.Warn("Service Warning", logging.ServiceLogging(
			"UserService", "RequestEmailUpdate",
			logging.DeleteEmailUpdateRequestError, err,
		)...)
	}

	token, err := security.GenerateSecureToken()
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "RequestEmailUpdate", logging.HasherError, err,
		)...)
		return ErrInternal
	}

	req := models.EmailUpdateRequest{
		UserID:    ID,
		OldEmail:  user.Email,
		NewEmail:  input.NewEmail,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.EmailUpdateRepo.CreateRequest(ctx, &req); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "RequestEmailUpdate", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	task, err := tasks.NewEmailUpdateConfirmationTask(user.Email, user.Name,
		user.Email, input.NewEmail, token, req.ExpiresAt)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "RequestEmailUpdate", logging.AsynqTaskError, err,
		)...)
		return ErrInternal
	}

	if _, err := s.AsynqClient.EnqueueContext(ctx, task,
		asynq.Queue(tasks.QueueEmail)); err != nil {
		if errors.Is(err, asynq.ErrDuplicateTask) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"UserService", "RequestEmailUpdate", logging.AsynqTaskError, err,
			)...)
			return ErrDuplicateTask
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "RequestEmailUpdate", logging.AsynqTaskError, err,
		)...)
		return ErrInternal
	}

	return nil
}

func (s *userService) ConfirmEmailUpdate(ctx context.Context, ID uuid.UUID,
	input models.ConfirmEmailUpdateInput) error {
	req, err := s.EmailUpdateRepo.GetByToken(ctx, input.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidEmailUpdateToken
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "ConfirmEmailUpdate", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if req.UserID != ID {
		return ErrInvalidEmailUpdateToken
	}

	if req.IsExpired() {
		return ErrExpiredEmailUpdateToken
	}

	user, err := s.Repo.GetUserByID(ctx, ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "ConfirmEmailUpdate", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	user.Email = req.NewEmail
	if err := s.Repo.UpdateUser(ctx, user); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"UserService", "ConfirmEmailUpdate", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if err := s.EmailUpdateRepo.DeleteRequestsByUserID(ctx, ID); err != nil {
		err = fmt.Errorf("%s: %v", ID, err)
		s.Logger.Warn("Service Warning", logging.ServiceLogging(
			"UserService", "ConfirmEmailUpdate",
			logging.DeleteEmailUpdateRequestError, err,
		)...)
	}

	return nil
}
