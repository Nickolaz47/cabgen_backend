package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	adminRouter := r.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())

	userGroup := adminRouter.Group("/user")
	userGroup.GET("", admin.GetAllUsers)
	userGroup.GET("/:username", admin.GetUserByUsername)
	userGroup.POST("", admin.CreateUser)
	userGroup.PUT("/:username")
	userGroup.DELETE("/:username")
}
