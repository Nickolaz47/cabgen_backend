package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/gin-gonic/gin"
)

func SetupOriginRoutes(r *gin.RouterGroup) {
	originRouter := r.Group("/origin")

	originRouter.GET("", origin.GetAllOrigins)
	originRouter.GET("/:originId", origin.GetOriginByID)
	originRouter.GET("/search", origin.GetOriginByName)
	originRouter.POST("", origin.CreateOrigin)
	originRouter.PUT("/:originId", origin.UpdateOrigin)
	originRouter.DELETE("/:originId", origin.DeleteOrigin)
}
