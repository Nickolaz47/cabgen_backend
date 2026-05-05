package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/healthservice"
	"github.com/gin-gonic/gin"
)

func SetupHealthServiceRoutes(
	r *gin.RouterGroup, handler *healthservice.HealthServiceHandler) {
	healthServiceRouter := r.Group("/health-services")
	healthServiceRouter.GET("", handler.GetActiveHealthServices)
}
