package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/origin"
	"github.com/gin-gonic/gin"
)

func SetupOriginRoutes(r *gin.RouterGroup) {
	originRouter := r.Group("/origin")
	originRouter.GET("", origin.GetActiveOrigins)
}
