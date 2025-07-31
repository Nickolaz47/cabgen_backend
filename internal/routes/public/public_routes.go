package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/gin-gonic/gin"
)

func PublicRoutes(r *gin.RouterGroup) {
	r.GET("/health", public.Health)

	authRouter := r.Group("/auth")
	authRouter.POST("/register", public.Register)
	authRouter.POST("/login", public.Login)
	authRouter.POST("/logout", public.Logout)
	authRouter.GET("/refresh", public.Refresh)
}
