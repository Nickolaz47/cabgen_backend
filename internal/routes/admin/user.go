package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/gin-gonic/gin"
)

func SetupAdminUserRoutes(r *gin.RouterGroup, handler *user.AdminUserHandler) {
	userRouter := r.Group("/users")

	userRouter.GET("", handler.GetUsers)
	userRouter.GET("/:userId", handler.GetUserByID)
	userRouter.POST("", handler.CreateUser)
	userRouter.PUT("/:userId", handler.UpdateUser)
	userRouter.PATCH("/:userId/activate", handler.ActivateUser)
	userRouter.PATCH("/:userId/deactivate", handler.DeactivateUser)
	userRouter.DELETE("/:userId", handler.DeleteUser)
}
