package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/microorganism"
	"github.com/gin-gonic/gin"
)

func SetupAdminMicroorganismRoutes(r *gin.RouterGroup, handler *microorganism.AdminMicroorganismHandler) {
	microorganismRouter := r.Group("/microorganisms")

	microorganismRouter.GET("", handler.GetMicroorganisms)
	microorganismRouter.GET("/:microorganismId", handler.GetMicroorganismByID)
	microorganismRouter.GET("/search", handler.GetMicroorganismBySpecies)
	microorganismRouter.POST("", handler.CreateMicroorganism)
	microorganismRouter.PUT("/:microorganismId", handler.UpdateMicroorganism)
	microorganismRouter.DELETE("/:microorganismId", handler.DeleteMicroorganism)
}
