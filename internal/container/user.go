package container

import (
	adminUser "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildUserService(db *gorm.DB, logger *zap.Logger) services.UserService {
	userRepo := repositories.NewUserRepo(db)
	countryRepo := repositories.NewCountryRepo(db)
	userService := services.NewUserService(userRepo, countryRepo, logger)

	return userService
}

func BuildAdminUserService(db *gorm.DB, logger *zap.Logger) services.AdminUserService {
	userRepo := repositories.NewUserRepo(db)
	countryRepo := repositories.NewCountryRepo(db)
	hasher := security.NewPasswordHasher()
	adminUserService := services.NewAdminUserService(
		userRepo, countryRepo, hasher, logger)

	return adminUserService
}

func BuildUserHandler(svc services.UserService) *user.UserHandler {
	return user.NewUserHandler(svc)
}

func BuildAdminUserHandler(svc services.AdminUserService) *adminUser.AdminUserHandler {
	return adminUser.NewAdminUserHandler(svc)
}
