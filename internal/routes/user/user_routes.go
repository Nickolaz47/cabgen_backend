package user

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	userRouter := r.Group("/user", middlewares.AuthMiddleware())
	userRouter.GET("/me", user.GetOwnUser)
	userRouter.PUT("/me", user.UpdateUser)
}
