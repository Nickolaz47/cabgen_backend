package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/gin-gonic/gin"
)

func SetupAdminSampleSourceRoutes(r *gin.RouterGroup, handler *samplesource.AdminSampleSourceHandler) {
	sampleSourceRouter := r.Group("/sample-sources")

	sampleSourceRouter.GET("", handler.GetSampleSources)
	sampleSourceRouter.GET("/:sampleSourceId", handler.GetSampleSourceByID)
	sampleSourceRouter.GET("/search", handler.GetSampleSourcesByNameOrGroup)
	sampleSourceRouter.POST("", handler.CreateSampleSource)
	sampleSourceRouter.PUT("/:sampleSourceId", handler.UpdateSampleSource)
	sampleSourceRouter.DELETE("/:sampleSourceId", handler.DeleteSampleSource)
}
