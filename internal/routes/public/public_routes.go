package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/gin-gonic/gin"
)

func PublicRoutes(r *gin.Engine) {
	r.GET("/health", public.Health)
	r.POST("/register", public.Register)
	r.POST("/login", public.Login)
	r.POST("/logout", public.Logout)
	r.GET("/refresh", public.Refresh)
}
