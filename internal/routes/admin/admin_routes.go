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
	userGroup.PUT("/:username", admin.UpdateUser)
	userGroup.PUT("/activation/:username", admin.UpdateUserActivation)
	userGroup.DELETE("/:username", admin.DeleteUser)

	originGroup := adminRouter.Group("/origin")
	originGroup.GET("", admin.GetAllOrigins)
	originGroup.GET("/:originId", admin.GetOriginByID)
	originGroup.GET("/search", admin.GetOriginByName)
	originGroup.POST("", admin.CreateOrigin)
	originGroup.PUT("/:originId", admin.UpdateOrigin)
	originGroup.DELETE("/:originId", admin.DeleteOrigin)
}
