package container

import (
	adminUser "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildUserService(db *gorm.DB) services.UserService {
	userRepo := repository.NewUserRepo(db)
	countryRepo := repository.NewCountryRepo(db)
	userService := services.NewUserService(userRepo, countryRepo)

	return userService
}

func BuildAdminUserService(db *gorm.DB) services.AdminUserService {
	userRepo := repository.NewUserRepo(db)
	countryRepo := repository.NewCountryRepo(db)
	hasher := security.NewPasswordHasher()
	adminUserService := services.NewAdminUserService(
		userRepo, countryRepo, hasher)

	return adminUserService
}

func BuildUserHandler(svc services.UserService) *user.UserHandler {
	return user.NewUserHandler(svc)
}

func BuildAdminUserHandler(svc services.AdminUserService) *adminUser.AdminUserHandler {
	return adminUser.NewAdminUserHandler(svc)
}
