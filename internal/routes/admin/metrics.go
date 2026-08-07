package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/metrics"
	"github.com/gin-gonic/gin"
)

func SetupAdminMetricsRoutes(r *gin.RouterGroup, handler *metrics.AdminMetricsHandler) {
	r.GET("/metrics", handler.GetMetrics)
}
