package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupPublicAuthRoutes(r *gin.RouterGroup, handler *auth.AuthHandler) {
	authRouter := r.Group("/auth")

	authRouter.POST("/register",
		middlewares.RateLimitPerMinute(3), handler.Register)
	authRouter.POST("/login",
		middlewares.RateLimitPerMinute(5), handler.Login)
	authRouter.POST("/refresh", handler.Refresh)
	authRouter.POST("/forgot-password",
		middlewares.RateLimitPerMinute(3), handler.ForgotPassword)
	authRouter.POST("/reset-password",
		middlewares.RateLimitPerMinute(3), handler.ResetPassword)
	authRouter.POST("/logout", handler.Logout)
}
