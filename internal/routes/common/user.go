package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, handler *user.UserHandler) {
	userRouter := r.Group("/users")

	userRouter.GET("/me", handler.GetOwnUser)
	userRouter.PUT("/me", handler.UpdateUser)
}
