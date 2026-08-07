package container

import (
	adminHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/metrics"
	publicHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/metrics"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildMetricsService(db *gorm.DB, logger *zap.Logger) services.MetricsService {
	sampleRepo := repositories.NewSampleRepo(db)
	analysisRepo := repositories.NewAnalysisRepository(db)
	userRepo := repositories.NewUserRepo(db)
	return services.NewMetricsService(sampleRepo, analysisRepo, userRepo, logger)
}

func BuildMetricsHandler(svc services.MetricsService) *publicHandler.MetricsHandler {
	return publicHandler.NewMetricsHandler(svc)
}

func BuildAdminMetricsHandler(svc services.MetricsService) *adminHandler.AdminMetricsHandler {
	return adminHandler.NewAdminMetricsHandler(svc)
}
