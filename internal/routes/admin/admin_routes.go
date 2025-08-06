package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	adminRouter := r.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())

	userGroup := adminRouter.Group("/user")
	userGroup.GET("")
	userGroup.GET("/:id")
	userGroup.POST("")
	userGroup.PUT("/:id")
	userGroup.DELETE("/:id")
}
