package services

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
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
}

type userService struct {
	Repo        repositories.UserRepository
	CountryRepo repositories.CountryRepository
	AsynqClient TaskEnqueuer
	Logger      *zap.Logger
	RootDir     string
}

func NewUserService(
	repo repositories.UserRepository,
	countryRepo repositories.CountryRepository,
	asynqClient TaskEnqueuer,
	logger *zap.Logger,
	rootDir string,
) UserService {
	return &userService{
		Repo:        repo,
		CountryRepo: countryRepo,
		AsynqClient: asynqClient,
		Logger:      logger,
		RootDir:     rootDir,
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
