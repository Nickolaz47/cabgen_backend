package services

import (
	"context"
	"errors"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminUserService interface {
	Find(ctx context.Context, filter models.AdminUserFilter, language string) ([]models.AdminUserResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID, language string) (*models.AdminUserResponse, error)
	FindByUsername(ctx context.Context, username, language string) (*models.AdminUserResponse, error)
	FindByEmail(ctx context.Context, email, language string) (*models.AdminUserResponse, error)
	Create(ctx context.Context, input models.AdminUserCreateInput, adminName, language string) (*models.AdminUserResponse, error)
	Update(ctx context.Context, ID uuid.UUID, input models.AdminUserUpdateInput, language string) (*models.AdminUserResponse, error)
	ActivateUser(ctx context.Context, ID uuid.UUID, adminName string) error
	DeactivateUser(ctx context.Context, ID uuid.UUID) error
	Delete(ctx context.Context, ID uuid.UUID) error
}

type adminUserService struct {
	Repo        repositories.UserRepository
	CountryRepo repositories.CountryRepository
	Hasher      security.PasswordHasher
	Logger      *zap.Logger
}

func NewAdminUserService(
	repo repositories.UserRepository,
	countryRepo repositories.CountryRepository,
	hasher security.PasswordHasher,
	logger *zap.Logger,
) AdminUserService {
	return &adminUserService{
		Repo:        repo,
		CountryRepo: countryRepo,
		Hasher:      hasher,
		Logger:      logger,
	}
}

func (s *adminUserService) Find(
	ctx context.Context,
	filter models.AdminUserFilter,
	language string) ([]models.AdminUserResponse, error) {
	users, err := s.Repo.GetUsers(ctx, filter)
	if err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Find", logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	responses := make([]models.AdminUserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToAdminResponse(language)
	}

	return responses, nil
}

func (s *adminUserService) FindByID(ctx context.Context, ID uuid.UUID, language string) (*models.AdminUserResponse, error) {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "FindByID", logging.DatabaseNotFoundError, err,
			)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "FindByID", logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	response := user.ToAdminResponse(language)
	return &response, nil
}

func (s *adminUserService) FindByUsername(ctx context.Context, username, language string) (*models.AdminUserResponse, error) {
	user, err := s.Repo.GetUserByUsername(ctx, username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "FindByUsername", logging.DatabaseNotFoundError, err,
			)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "FindByUsername", logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	response := user.ToAdminResponse(language)
	return &response, nil
}

func (s *adminUserService) FindByEmail(ctx context.Context, email, language string) (*models.AdminUserResponse, error) {
	user, err := s.Repo.GetUserByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "FindByEmail", logging.DatabaseNotFoundError, err,
			)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "FindByEmail", logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	response := user.ToAdminResponse(language)
	return &response, nil
}

func (s *adminUserService) Create(
	ctx context.Context,
	input models.AdminUserCreateInput,
	adminName, language string) (*models.AdminUserResponse, error) {
	existingUser, err := s.Repo.ExistsByEmail(ctx, &input.Email, uuid.Nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Create", logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}
	if existingUser != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Create",
				logging.DatabaseConflictEmailError, err,
			)...)
		return nil, ErrConflictEmail
	}

	existingUser, err = s.Repo.ExistsByUsername(ctx, &input.Username, uuid.Nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Create",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}
	if existingUser != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Create",
				logging.DatabaseConflictUsernameError, err,
			)...)
		return nil, ErrConflictUsername
	}

	hashedPassword, err := s.Hasher.Hash(input.Password)
	if err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Create",
				logging.HasherError, err,
			)...)
		return nil, ErrInternal
	}

	country, err := s.CountryRepo.GetCountryByCode(ctx, input.CountryCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "Create",
					logging.ExternalDatabaseNotFoundError, err,
				)...)
			return nil, ErrInvalidCountryCode
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Create",
				logging.ExternalDatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	activatedOn := time.Now()
	user := models.User{
		Name:        input.Name,
		Username:    input.Username,
		Email:       input.Email,
		Password:    hashedPassword,
		CountryID:   country.ID,
		UserRole:    input.UserRole,
		IsActive:    input.IsActive,
		Interest:    input.Interest,
		Role:        input.Role,
		Institution: input.Institution,
		CreatedBy:   adminName,
		ActivatedBy: &adminName,
		ActivatedOn: &activatedOn,
	}

	if err := s.Repo.CreateUser(ctx, &user); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Create",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	user.Country = *country

	response := user.ToAdminResponse(language)
	return &response, nil
}

func (s *adminUserService) Update(
	ctx context.Context, ID uuid.UUID,
	input models.AdminUserUpdateInput,
	language string) (*models.AdminUserResponse, error) {
	existingUser, err := s.Repo.GetUserByID(
		ctx, ID,
	)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Update",
				logging.DatabaseNotFoundError, err,
			)...)
		return nil, ErrNotFound
	}
	if err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Update",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	if input.Email != nil {
		duplicate, err := s.Repo.ExistsByEmail(ctx, input.Email, ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "Update",
					logging.DatabaseError, err,
				)...)
			return nil, ErrInternal
		}
		if duplicate != nil {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "Update",
					logging.DatabaseConflictEmailError, err,
				)...)
			return nil, ErrConflictEmail
		}
	}

	if input.Username != nil {
		duplicate, err := s.Repo.ExistsByUsername(ctx, input.Username, ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "Update",
					logging.DatabaseError, err,
				)...)
			return nil, ErrInternal
		}
		if duplicate != nil {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "Update",
					logging.DatabaseConflictUsernameError, err,
				)...)
			return nil, ErrConflictUsername
		}
	}

	var hashedPassword string
	if input.Password != nil {
		hashedPassword, err = s.Hasher.Hash(*input.Password)
		if err != nil {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "Update",
					logging.HasherError, err,
				)...)
			return nil, ErrInternal
		}

		input.Password = &hashedPassword
	}

	if input.CountryCode != nil {
		country, err := s.CountryRepo.GetCountryByCode(ctx, *input.CountryCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"AdminUserService", "Update",
						logging.ExternalDatabaseNotFoundError, err,
					)...)
				return nil, ErrInvalidCountryCode
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "Update",
					logging.ExternalDatabaseError, err,
				)...)
			return nil, ErrInternal
		}
		existingUser.CountryID = country.ID
		existingUser.Country = *country
	}

	validations.ApplyAdminUpdateToUser(existingUser, &input)

	if err := s.Repo.UpdateUser(ctx, existingUser); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Update",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	response := existingUser.ToAdminResponse(language)
	return &response, nil
}

func (s *adminUserService) ActivateUser(ctx context.Context, ID uuid.UUID, adminName string) error {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "ActivateUser",
					logging.DatabaseNotFoundError, err,
				)...)
			return ErrNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "ActivateUser",
				logging.DatabaseError, err,
			)...)
		return ErrInternal
	}

	if user.IsActive {
		return nil
	}

	activatedOn := time.Now()
	user.IsActive = true
	user.ActivatedBy = &adminName
	user.ActivatedOn = &activatedOn

	if err := s.Repo.UpdateUser(ctx, user); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "ActivateUser",
				logging.DatabaseError, err,
			)...)
		return ErrInternal
	}

	return nil
}

func (s *adminUserService) DeactivateUser(ctx context.Context, ID uuid.UUID) error {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminUserService", "DeactivateUser",
					logging.DatabaseNotFoundError, err,
				)...)
			return ErrNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "DeactivateUser",
				logging.DatabaseError, err,
			)...)
		return ErrInternal
	}

	if !user.IsActive {
		return nil
	}

	user.IsActive = false
	user.ActivatedBy = nil
	user.ActivatedOn = nil

	if err := s.Repo.UpdateUser(ctx, user); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "DeactivateUser",
				logging.DatabaseError, err,
			)...)
		return ErrInternal
	}

	return nil
}

func (s *adminUserService) Delete(ctx context.Context, ID uuid.UUID) error {
	user, err := s.Repo.GetUserByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Delete",
				logging.DatabaseNotFoundError, err,
			)...)
		return ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Delete",
				logging.DatabaseError, err,
			)...)
		return ErrInternal
	}

	if err := s.Repo.DeleteUser(ctx, user); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminUserService", "Delete",
				logging.DatabaseError, err,
			)...)
		return ErrInternal
	}

	return nil
}
