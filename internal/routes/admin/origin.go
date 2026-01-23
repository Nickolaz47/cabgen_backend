package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/gin-gonic/gin"
)

func SetupAdminOriginRoutes(r *gin.RouterGroup, handler *origin.AdminOriginHandler) {
	originRouter := r.Group("/origins")

	originRouter.GET("", handler.GetAllOrigins)
	originRouter.GET("/:originId", handler.GetOriginByID)
	originRouter.GET("/search", handler.GetOriginsByName)
	originRouter.POST("", handler.CreateOrigin)
	originRouter.PUT("/:originId", handler.UpdateOrigin)
	originRouter.DELETE("/:originId", handler.DeleteOrigin)
}
