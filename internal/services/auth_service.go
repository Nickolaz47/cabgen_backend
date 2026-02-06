package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/events"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, input models.UserRegisterInput, language string) (*models.UserResponse, error)
	Login(ctx context.Context, input models.LoginInput) (*models.Cookies, error)
	Refresh(ctx context.Context, tokenStr string) (*http.Cookie, error)
}

type authService struct {
	UserRepo      repositories.UserRepository
	CountryRepo   repositories.CountryRepository
	EventEmitter  events.EventEmitter
	Hasher        security.PasswordHasher
	TokenProvider auth.TokenProvider
}

func NewAuthService(
	userRepo repositories.UserRepository,
	countryRepo repositories.CountryRepository,
	emitter events.EventEmitter,
	hasher security.PasswordHasher,
	tokenProvider auth.TokenProvider,
) AuthService {
	return &authService{
		UserRepo:      userRepo,
		CountryRepo:   countryRepo,
		EventEmitter:  emitter,
		Hasher:        hasher,
		TokenProvider: tokenProvider,
	}
}

func (s *authService) Register(
	ctx context.Context,
	input models.UserRegisterInput,
	language string,
) (*models.UserResponse, error) {
	existingUser, err := s.UserRepo.ExistsByEmail(
		ctx, &input.Email, uuid.Nil,
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}
	if existingUser != nil {
		return nil, ErrConflictEmail
	}

	existingUser, err = s.UserRepo.ExistsByUsername(ctx, &input.Username, uuid.Nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}
	if existingUser != nil {
		return nil, ErrConflictUsername
	}

	if ok := validations.IsEmailMatch(
		input.Email, input.ConfirmEmail); !ok {
		return nil, ErrEmailMismatch
	}

	if ok := validations.IsPasswordMatch(
		input.Password, input.ConfirmPassword,
	); !ok {
		return nil, ErrPasswordMismatch
	}

	hashedPassword, err := s.Hasher.Hash(input.Password)
	if err != nil {
		return nil, ErrInternal
	}

	country, err := s.CountryRepo.GetCountryByCode(ctx, input.CountryCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCountryCode
		}
		return nil, ErrInternal
	}

	user := models.User{
		Name:        input.Name,
		Username:    input.Username,
		Email:       input.Email,
		Password:    hashedPassword,
		CountryID:   country.ID,
		UserRole:    models.Collaborator,
		Interest:    input.Interest,
		Role:        input.Role,
		Institution: input.Institution,
		CreatedBy:   input.Username,
	}

	if err := s.UserRepo.CreateUser(ctx, &user); err != nil {
		return nil, ErrInternal
	}

	user.Country = *country

	s.EventEmitter.Emit(ctx, events.EventUserRegistered,
		events.UserRegisteredPayload{RegisteredUsername: user.Username})

	response := user.ToResponse(language)
	return &response, nil
}

func (s *authService) Login(
	ctx context.Context, input models.LoginInput) (*models.Cookies, error) {
	existingUser, err := s.UserRepo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, ErrInternal
	}

	if !existingUser.IsActive {
		return nil, ErrDisabledUser
	}

	if err = s.Hasher.CheckPassword(existingUser.Password,
		input.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, ErrInternal
	}

	accessKey, err := auth.GetSecretKey(auth.Access)
	if err != nil {
		return nil, ErrInternal
	}

	refreshKey, err := auth.GetSecretKey(auth.Refresh)
	if err != nil {
		return nil, ErrInternal
	}

	accessToken, err := s.TokenProvider.GenerateToken(
		existingUser.ToToken(), accessKey, auth.AccessTokenExpiration)
	if err != nil {
		return nil, ErrInternal
	}

	refreshToken, err := s.TokenProvider.GenerateToken(
		existingUser.ToToken(), refreshKey, auth.RefreshTokenExpiration)
	if err != nil {
		return nil, ErrInternal
	}

	return &models.Cookies{
		AccessCookie: auth.CreateCookie(
			auth.Access, accessToken, "/",
			auth.AccessTokenExpiration,
		),
		RefreshCookie: auth.CreateCookie(
			auth.Refresh, refreshToken,
			"/api/auth/refresh", auth.RefreshTokenExpiration,
		),
	}, nil
}

func (s *authService) Refresh(ctx context.Context, tokenStr string) (*http.Cookie, error) {
	refreshSecret, err := auth.GetSecretKey(auth.Refresh)
	if err != nil {
		return nil, ErrInternal
	}

	userToken, err := s.TokenProvider.ValidateToken(tokenStr, refreshSecret)
	if err != nil {
		return nil, ErrUnauthorized
	}

	accessSecret, err := auth.GetSecretKey(auth.Access)
	if err != nil {
		return nil, ErrInternal
	}

	accessToken, err := s.TokenProvider.GenerateToken(
		*userToken, accessSecret, auth.AccessTokenExpiration)
	if err != nil {
		return nil, ErrInternal
	}

	return auth.CreateCookie(
		auth.Access, accessToken,
		"/", auth.AccessTokenExpiration,
	), nil
}
