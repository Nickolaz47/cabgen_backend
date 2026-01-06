package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/origin"
	"github.com/gin-gonic/gin"
)

func SetupOriginRoutes(r *gin.RouterGroup, handler *origin.OriginHandler) {
	originRouter := r.Group("/origin")
	originRouter.GET("", handler.GetActiveOrigins)
}
