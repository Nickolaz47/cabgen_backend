package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/samplesource"
	"github.com/gin-gonic/gin"
)

func SetupSampleSourceRoutes(r *gin.RouterGroup, handler *samplesource.SampleSourceHandler) {
	sampleSourceRouter := r.Group("/sample-sources")
	sampleSourceRouter.GET("", handler.GetActiveSampleSources)
}
