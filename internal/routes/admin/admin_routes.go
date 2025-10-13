package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	adminRouter := r.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())

	userGroup := adminRouter.Group("/user")
	userGroup.GET("", user.GetAllUsers)
	userGroup.GET("/:username", user.GetUserByUsername)
	userGroup.POST("", user.CreateUser)
	userGroup.PUT("/:username", user.UpdateUser)
	userGroup.PUT("/activation/:username", user.UpdateUserActivation)
	userGroup.DELETE("/:username", user.DeleteUser)

	originGroup := adminRouter.Group("/origin")
	originGroup.GET("", origin.GetAllOrigins)
	originGroup.GET("/:originId", origin.GetOriginByID)
	originGroup.GET("/search", origin.GetOriginByName)
	originGroup.POST("", origin.CreateOrigin)
	originGroup.PUT("/:originId", origin.UpdateOrigin)
	originGroup.DELETE("/:originId", origin.DeleteOrigin)

	sequencerGroup := adminRouter.Group("/sequencer")
	sequencerGroup.GET("", sequencer.GetAllSequencers)
	sequencerGroup.GET("/:sequencerId", sequencer.GetSequencerByID)
	sequencerGroup.GET("/search", sequencer.GetSequencersByBrandOrModel)
	sequencerGroup.POST("", sequencer.CreateSequencer)
	sequencerGroup.PUT("/:sequencerId", sequencer.UpdateSequencer)
	sequencerGroup.DELETE("/:sequencerId", sequencer.DeleteSequencer)
}
