package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/analysis"
	"github.com/gin-gonic/gin"
)

func SetupAnalysisRoutes(r *gin.RouterGroup,
	handler *analysis.AnalysisHandler) {
	analysisRouter := r.Group("/analyses")

	analysisRouter.GET("", handler.GetAnalyses)
	analysisRouter.GET("/:analysisId", handler.GetAnalysisByID)
	analysisRouter.GET("/:analysisId/download/tsv", handler.DownloadZip)
	analysisRouter.POST("", handler.CreateAnalysis)
	analysisRouter.POST("/download/tsv", handler.DownloadBatchTSV)
	analysisRouter.DELETE("/:analysisId", handler.DeleteAnalysis)
}
