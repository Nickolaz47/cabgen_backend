package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildAnalysisRunnerService(db *gorm.DB, config pipeline.ToolsConfig,
	asynqClient *asynq.Client, rootDir string,
	logger *zap.Logger) services.AnalysisRunnerService {
	analysisRepo := repositories.NewAnalysisRepository(db)
	runner := pipeline.NewToolRunner(&pipeline.RealCommander{})
	pipeline := pipeline.NewCabgenPipeline(runner, config)

	return services.NewAnalysisRunnerService(
		analysisRepo, pipeline, asynqClient, logger, rootDir,
	)
}
