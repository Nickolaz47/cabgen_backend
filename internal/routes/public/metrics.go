package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/metrics"
	"github.com/gin-gonic/gin"
)

func SetupMetricsRoutes(r *gin.RouterGroup, handler *metrics.MetricsHandler) {
	r.GET("/metrics", handler.GetMetrics)
}
