package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/metrics"
	"github.com/gin-gonic/gin"
)

func SetupUserMetricsRoutes(r *gin.RouterGroup,
	handler *metrics.MetricsHandler) {
	userRouter := r.Group("/users")

	userRouter.GET("/me/metrics", handler.GetUserMetrics)
}
