package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/gin-gonic/gin"
)

func SetupSequencerRoutes(r *gin.RouterGroup) {
	sequencerRouter := r.Group("/sequencer")

	sequencerRouter.GET("", sequencer.GetAllSequencers)
	sequencerRouter.GET("/:sequencerId", sequencer.GetSequencerByID)
	sequencerRouter.GET("/search", sequencer.GetSequencersByBrandOrModel)
	sequencerRouter.POST("", sequencer.CreateSequencer)
	sequencerRouter.PUT("/:sequencerId", sequencer.UpdateSequencer)
	sequencerRouter.DELETE("/:sequencerId", sequencer.DeleteSequencer)
}
