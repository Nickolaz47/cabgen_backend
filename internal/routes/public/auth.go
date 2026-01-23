package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup, handler *auth.AuthHandler) {
	authRouter := r.Group("/auth")

	authRouter.POST("/register", handler.Register)
	authRouter.POST("/login", handler.Login)
	authRouter.POST("/logout", handler.Logout)
	authRouter.POST("/refresh", handler.Refresh)
}
