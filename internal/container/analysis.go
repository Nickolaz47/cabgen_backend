package container

import (
	adminHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildAnalysisService(db *gorm.DB, asynqClient *asynq.Client,
	inspector *asynq.Inspector, logger *zap.Logger,
	rootDir string) services.AnalysisService {
	analysisRepo := repositories.NewAnalysisRepository(db)
	sampleRepo := repositories.NewSampleRepo(db)
	userRepo := repositories.NewUserRepo(db)
	analysisService := services.NewAnalysisService(
		analysisRepo, sampleRepo,
		userRepo, asynqClient, inspector, logger, rootDir,
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
