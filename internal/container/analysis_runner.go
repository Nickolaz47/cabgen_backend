package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildAnalysisRunnerService(db *gorm.DB, config pipeline.ToolsConfig,
	taskQueue services.TaskQueue, rootDir string,
	logger *zap.Logger) services.AnalysisRunnerService {
	analysisRepo := repositories.NewAnalysisRepository(db)
	cmdr := &pipeline.RealCommander{}
	runner := pipeline.NewToolRunner(cmdr)
	pipeline := pipeline.NewCabgenPipeline(runner, config, logger)

	return services.NewAnalysisRunnerService(
		analysisRepo, pipeline, cmdr, taskQueue, logger, rootDir,
	)
}
