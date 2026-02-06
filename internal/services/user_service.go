package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	FindByID(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.UserUpdateInput, language string) (*models.UserResponse, error)
}

type userService struct {
	Repo        repositories.UserRepository
	CountryRepo repositories.CountryRepository
}

func NewUserService(repo repositories.UserRepository, countryRepo repositories.CountryRepository) UserService {
	return &userService{Repo: repo, CountryRepo: countryRepo}
}

func (s *userService) FindByID(
	ctx context.Context, ID uuid.UUID,
	language string) (*models.UserResponse, error) {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
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
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternal
	}

	if input.Username != nil {
		duplicate, err := s.Repo.ExistsByUsername(
			ctx,
			input.Username,
			ID,
		)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInternal
		}
		if duplicate != nil {
			return nil, ErrConflictUsername
		}
	}

	if input.CountryCode != nil {
		country, err := s.CountryRepo.GetCountryByCode(ctx, *input.CountryCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrInvalidCountryCode
			}
			return nil, ErrInternal
		}
		existingUser.CountryID = country.ID
		existingUser.Country = *country
	}

	validations.ApplyUpdateToUser(existingUser, &input)

	if err := s.Repo.UpdateUser(ctx, existingUser); err != nil {
		return nil, ErrInternal
	}

	response := existingUser.ToResponse(language)
	return &response, nil
}
