package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	commonAuthHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/common/auth"
	publicAuthHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildAuthService(mainDB *gorm.DB,
	taskQueue services.TaskQueue, logger *zap.Logger) services.AuthService {
	countryRepo := repositories.NewCountryRepo(mainDB)
	userRepo := repositories.NewUserRepo(mainDB)
	passwordResetRepo := repositories.NewPasswordResetRepo(mainDB)
	hasher := security.NewPasswordHasher()
	provider := auth.NewTokenProvider()
	authService := services.NewAuthService(
		userRepo, countryRepo, passwordResetRepo, hasher, provider,
		taskQueue, logger,
	)

	return authService
}

func BuildPublicAuthHandler(
	svc services.AuthService) *publicAuthHandler.AuthHandler {
	return publicAuthHandler.NewAuthHandler(svc)
}

func BuildCommonAuthHandler(
	svc services.AuthService) *commonAuthHandler.AuthHandler {
	return commonAuthHandler.NewAuthHandler(svc)
}
