package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	authHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildAuthService(mainDB *gorm.DB, logger *zap.Logger) services.AuthService {
	countryRepo := repositories.NewCountryRepo(mainDB)
	userRepo := repositories.NewUserRepo(mainDB)
	hasher := security.NewPasswordHasher()
	provider := auth.NewTokenProvider()
	authService := services.NewAuthService(
		userRepo, countryRepo, hasher, provider, logger,
	)

	return authService
}

func BuildAuthHandler(svc services.AuthService) *authHandler.AuthHandler {
	return authHandler.NewAuthHandler(svc)
}
