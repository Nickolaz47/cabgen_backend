package user

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	r.GET("/user/me", middlewares.AuthMiddleware(), user.GetOwnUser)
	r.PATCH("/user/me", middlewares.AuthMiddleware(), user.UpdateUser)
}
