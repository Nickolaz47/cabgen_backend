package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/sequencer"
	"github.com/gin-gonic/gin"
)

func SetupSequencerRoutes(r *gin.RouterGroup) {
	sequencerRouter := r.Group("/sequencer")
	sequencerRouter.GET("", sequencer.GetActiveSequencers)
}
