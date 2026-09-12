package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupPublicAuthRoutes(r *gin.RouterGroup, handler *auth.AuthHandler) {
	authRouter := r.Group("/auth")

	authRouter.POST("/register", handler.Register)
	authRouter.POST("/login", handler.Login)
	authRouter.POST("/refresh", handler.Refresh)
	authRouter.POST("/forgot-password",
		middlewares.GlobalRateLimitPerMinute(3), handler.ForgotPassword)
	authRouter.POST("/reset-password",
		middlewares.GlobalRateLimitPerMinute(3), handler.ResetPassword)
}
