package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/healthservice"
	"github.com/gin-gonic/gin"
)

func SetupHealthServiceRoutes(
	r *gin.RouterGroup, handler *healthservice.HealthServiceHandler) {
	sampleSourceRouter := r.Group("/health-service")
	sampleSourceRouter.GET("", handler.GetActiveHealthServices)
}
