package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/gin-gonic/gin"
)

func SetupAdminAnalysisRoutes(r *gin.RouterGroup,
	handler *analysis.AdminAnalysisHandler) {
	analysisRouter := r.Group("/analyses")

	analysisRouter.GET("", handler.GetAnalyses)
	analysisRouter.GET("/:analysisId", handler.GetAnalysisByID)
	analysisRouter.GET("/:analysisId/download/tsv", handler.DownloadZip)
	analysisRouter.GET("/types", handler.GetAnalysisTypes)
	analysisRouter.POST("", handler.CreateAnalysis)
	analysisRouter.POST("/download/tsv", handler.DownloadBatchTSV)
	analysisRouter.PUT("/:analysisId", handler.UpdateAnalysis)
	analysisRouter.DELETE("/:analysisId", handler.DeleteAnalysis)
}
