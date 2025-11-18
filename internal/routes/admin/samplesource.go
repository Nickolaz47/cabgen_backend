package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/gin-gonic/gin"
)

func SetupSampleSourceRoutes(r *gin.RouterGroup) {
	sampleSourceRouter := r.Group("/sampleSource")

	sampleSourceRouter.GET("", samplesource.GetSampleSources)
	sampleSourceRouter.GET("/:sampleSourceId", samplesource.GetSampleSourceByID)
	sampleSourceRouter.GET("/search", samplesource.GetSampleSourceByNameOrGroup)
	sampleSourceRouter.POST("", samplesource.CreateSampleSource)
	sampleSourceRouter.PUT("/:sampleSourceId", samplesource.UpdateSampleSource)
	sampleSourceRouter.DELETE("/:sampleSourceId", samplesource.DeleteSampleSource)
}
