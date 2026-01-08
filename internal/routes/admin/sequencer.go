package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/gin-gonic/gin"
)

func SetupAdminSequencerRoutes(r *gin.RouterGroup, handler *sequencer.AdminSequencerHandler) {
	sequencerRouter := r.Group("/sequencer")

	sequencerRouter.GET("", handler.GetSequencers)
	sequencerRouter.GET("/:sequencerId", handler.GetSequencerByID)
	sequencerRouter.GET("/search", handler.GetSequencersByBrandOrModel)
	sequencerRouter.POST("", handler.CreateSequencer)
	sequencerRouter.PUT("/:sequencerId", handler.UpdateSequencer)
	sequencerRouter.DELETE("/:sequencerId", handler.DeleteSequencer)
}
