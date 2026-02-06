package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/events"
	authHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildAuthService(mainDB, eventDB *gorm.DB) services.AuthService {
	countryRepo := repositories.NewCountryRepo(mainDB)
	userRepo := repositories.NewUserRepo(mainDB)
	emitter := events.NewEventEmitter(repositories.NewEventRepo(eventDB))
	hasher := security.NewPasswordHasher()
	provider := auth.NewTokenProvider()
	authService := services.NewAuthService(
		userRepo, countryRepo, emitter, hasher, provider,
	)

	return authService
}

func BuildAuthHandler(svc services.AuthService) *authHandler.AuthHandler {
	return authHandler.NewAuthHandler(svc)
}
