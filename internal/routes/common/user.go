package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup) {
	userRouter := r.Group("/user")

	userRouter.GET("/me", user.GetOwnUser)
	userRouter.PUT("/me", user.UpdateUser)
}
