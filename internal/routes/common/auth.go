package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/auth"
	"github.com/gin-gonic/gin"
)

func SetupCommonAuthRoutes(r *gin.RouterGroup, handler *auth.AuthHandler) {
	authRouter := r.Group("/auth")

	authRouter.POST("/logout", handler.Logout)
	authRouter.GET("/me", handler.Me)
}
