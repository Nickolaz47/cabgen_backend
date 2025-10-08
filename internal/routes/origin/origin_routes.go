package origin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func OriginRoutes(r *gin.RouterGroup) {
	originRouter := r.Group("/origin", middlewares.AuthMiddleware())
	originRouter.GET("", origin.GetAllActiveOrigins)
}
