package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/health"
	"github.com/gin-gonic/gin"
)

func SetupHealthRoute(r *gin.RouterGroup, handler *health.HealthHandler) {
	r.GET("/health", handler.Health)
}
