package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/microorganism"
	"github.com/gin-gonic/gin"
)

func SetupMicroorganismRoutes(
	r *gin.RouterGroup,handler *microorganism.MicroorganismHandler) {
	microorganismRouter := r.Group("/microorganisms")
	microorganismRouter.GET("", handler.GetActiveMicroorganisms)
}