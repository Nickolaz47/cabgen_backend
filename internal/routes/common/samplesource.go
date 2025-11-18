package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/samplesource"
	"github.com/gin-gonic/gin"
)

func SetupSampleSourceRoutes(r *gin.RouterGroup) {
	sampleSourceRouter := r.Group("/sampleSource")
	sampleSourceRouter.GET("", samplesource.GetActiveSampleSources)
}
