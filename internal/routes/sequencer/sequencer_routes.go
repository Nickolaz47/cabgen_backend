package sequencer

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SequencerRoutes(r *gin.RouterGroup) {
	sequencerRouter := r.Group("/sequencer", middlewares.AuthMiddleware())
	sequencerRouter.GET("", sequencer.GetActiveSequencers)
}
