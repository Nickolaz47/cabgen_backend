package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup) {
	authRouter := r.Group("/auth")

	authRouter.POST("/register", public.Register)
	authRouter.POST("/login", public.Login)
	authRouter.POST("/logout", public.Logout)
	authRouter.POST("/refresh", public.Refresh)
}
