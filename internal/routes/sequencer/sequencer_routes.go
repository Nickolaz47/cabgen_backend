package sequencer

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func OriginRoutes(r *gin.RouterGroup) {
	originRouter := r.Group("/sequencer", middlewares.AuthMiddleware())
	originRouter.GET("", sequencer.GetActiveSequencers)
}
