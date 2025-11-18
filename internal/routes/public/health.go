package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/gin-gonic/gin"
)

func SetupHealthRoute(r *gin.RouterGroup) {
	r.GET("/health", public.Health)
}
