package samplesource

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SampleSourceRoutes(r *gin.RouterGroup) {
	sampleSourceRouter := r.Group("/sampleSource", middlewares.AuthMiddleware())
	sampleSourceRouter.GET("", samplesource.GetActiveSampleSources)
}
