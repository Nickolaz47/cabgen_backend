package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/gin-gonic/gin"
)

func SetupPublicAuthRoutes(r *gin.RouterGroup, handler *auth.AuthHandler) {
	authRouter := r.Group("/auth")

	authRouter.POST("/register", handler.Register)
	authRouter.POST("/login", handler.Login)
	authRouter.POST("/refresh", handler.Refresh)
	authRouter.POST("/forgot-password", handler.ForgotPassword)
	authRouter.POST("/reset-password", handler.ResetPassword)
}
