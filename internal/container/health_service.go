package container

import (
	adminHealthService "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/healthservice"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/healthservice"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildHealthServiceService(
	db *gorm.DB, logger *zap.Logger) services.HealthServiceService {
	healthServiceRepo := repositories.NewHealthServiceRepo(db)
	countryRepo := repositories.NewCountryRepo(db)
	healthServiceService := services.NewHealthServiceService(
		healthServiceRepo, countryRepo, logger)

	return healthServiceService
}

func BuildHealthServiceHandler(
	svc services.HealthServiceService) *healthservice.HealthServiceHandler {
	return healthservice.NewHealthServiceHandler(svc)
}

func BuildAdminHealthServiceHandler(
	svc services.HealthServiceService) *adminHealthService.AdminHealthServiceHandler {
	return adminHealthService.NewAdminHealthServiceHandler(svc)
}
