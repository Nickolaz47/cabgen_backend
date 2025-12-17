package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sequencer"
	"github.com/gin-gonic/gin"
)

func SetupSequencerRoutes(r *gin.RouterGroup, handler *sequencer.SequencerHandler) {
	sequencerRouter := r.Group("/sequencer")
	sequencerRouter.GET("", handler.GetActiveSequencers)
}
