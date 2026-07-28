package container

import (
	adminUser "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/user"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildUserService(db *gorm.DB, asynqClient *asynq.Client,
	logger *zap.Logger, rootDir string) services.UserService {
	userRepo := repositories.NewUserRepo(db)
	countryRepo := repositories.NewCountryRepo(db)
	hasher := security.NewPasswordHasher()
	userService := services.NewUserService(
		userRepo, countryRepo, hasher, asynqClient, logger, rootDir)

	return userService
}

func BuildAdminUserService(db *gorm.DB, asynqClient *asynq.Client,
	logger *zap.Logger, rootDir string) services.AdminUserService {
	userRepo := repositories.NewUserRepo(db)
	countryRepo := repositories.NewCountryRepo(db)
	hasher := security.NewPasswordHasher()
	adminUserService := services.NewAdminUserService(
		userRepo, countryRepo, hasher, asynqClient, logger, rootDir)

	return adminUserService
}

func BuildUserHandler(svc services.UserService) *user.UserHandler {
	return user.NewUserHandler(svc)
}

func BuildAdminUserHandler(svc services.AdminUserService) *adminUser.AdminUserHandler {
	return adminUser.NewAdminUserHandler(svc)
}
