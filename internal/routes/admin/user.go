package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/gin-gonic/gin"
)

func SetupAdminUserRoutes(r *gin.RouterGroup) {
	userRouter := r.Group("/user")

	userRouter.GET("", user.GetAllUsers)
	userRouter.GET("/:username", user.GetUserByUsername)
	userRouter.POST("", user.CreateUser)
	userRouter.PUT("/:username", user.UpdateUser)
	userRouter.PUT("/activation/:username", user.UpdateUserActivation)
	userRouter.DELETE("/:username", user.DeleteUser)
}
