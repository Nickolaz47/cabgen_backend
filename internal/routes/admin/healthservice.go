package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/healthservice"
	"github.com/gin-gonic/gin"
)

func SetupAdminHealthServiceRoutes(
	r *gin.RouterGroup, handler *healthservice.AdminHealthServiceHandler) {
	healthServiceRouter := r.Group("/health-services")

	healthServiceRouter.GET("", handler.GetAllHealthServices)
	healthServiceRouter.GET("/:healthServiceId", handler.GetHealthServiceByID)
	healthServiceRouter.GET("/search", handler.GetHealthServicesByName)
	healthServiceRouter.POST("", handler.CreateHealthService)
	healthServiceRouter.PUT("/:healthServiceId", handler.UpdateHealthService)
	healthServiceRouter.DELETE("/:healthServiceId", handler.DeleteHealthService)
}
