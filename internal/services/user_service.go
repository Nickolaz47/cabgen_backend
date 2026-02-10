package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService interface {
	FindByID(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.UserUpdateInput, language string) (*models.UserResponse, error)
}

type userService struct {
	Repo        repositories.UserRepository
	CountryRepo repositories.CountryRepository
	Logger      *zap.Logger
}

func NewUserService(
	repo repositories.UserRepository,
	countryRepo repositories.CountryRepository,
	logger *zap.Logger,
) UserService {
	return &userService{Repo: repo, CountryRepo: countryRepo, Logger: logger}
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
