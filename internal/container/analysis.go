package container

import (
	adminHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildAnalysisService(db *gorm.DB,
	taskQueue services.TaskQueue, logger *zap.Logger,
	rootDir string) services.AnalysisService {
	analysisRepo := repositories.NewAnalysisRepository(db)
	sampleRepo := repositories.NewSampleRepo(db)
	userRepo := repositories.NewUserRepo(db)
	analysisService := services.NewAnalysisService(
		analysisRepo, sampleRepo,
		userRepo, taskQueue, logger, rootDir,
	)

	return analysisService
}

func BuildAnalysisHandler(svc services.AnalysisService,
) *analysis.AnalysisHandler {
	return analysis.NewAnalysisHandler(svc)
}

func BuildAdminAnalysisHandler(svc services.AnalysisService,
) *adminHandler.AdminAnalysisHandler {
	return adminHandler.NewAdminAnalysisHandler(svc)
}
