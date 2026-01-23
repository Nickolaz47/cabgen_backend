package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	authHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildAuthService(db *gorm.DB) services.AuthService {
	countryRepo := repository.NewCountryRepo(db)
	userRepo := repository.NewUserRepo(db)
	hasher := security.NewPasswordHasher()
	provider := auth.NewTokenProvider()
	authService := services.NewAuthService(
		userRepo, countryRepo, hasher, provider,
	)

	return authService
}

func BuildAuthHandler(svc services.AuthService) *authHandler.AuthHandler {
	return authHandler.NewAuthHandler(svc)
}
