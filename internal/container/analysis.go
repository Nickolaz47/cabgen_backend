package container

import (
	adminHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildAnalysisService(db *gorm.DB, logger *zap.Logger,
) services.AnalysisService {
	analysisRepo := repositories.NewAnalysisRepository(db)
	sampleRepo := repositories.NewSampleRepo(db)
	userRepo := repositories.NewUserRepo(db)
	analysisService := services.NewAnalysisService(
		analysisRepo, sampleRepo,
		userRepo, logger,
	)

	return analysisService
}

func BuildAdminAnalysisService(db *gorm.DB, logger *zap.Logger,
) services.AdminAnalysisService {
	analysisRepo := repositories.NewAnalysisRepository(db)
	sampleRepo := repositories.NewSampleRepo(db)
	userRepo := repositories.NewUserRepo(db)
	adminAnalysisService := services.NewAdminAnalysisService(
		analysisRepo, sampleRepo,
		userRepo, logger,
	)

	return adminAnalysisService
}

func BuildAnalysisHandler(svc services.AnalysisService,
) *analysis.AnalysisHandler {
	return analysis.NewAnalysisHandler(svc)
}

func BuildAdminAnalysisHandler(svc services.AdminAnalysisService,
) *adminHandler.AdminAnalysisHandler {
	return adminHandler.NewAdminAnalysisHandler(svc)
}
