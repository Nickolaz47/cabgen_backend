package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/selectoptions"
	"github.com/gin-gonic/gin"
)

func SetupSelectOptionRoutes(r *gin.RouterGroup,
	handler *selectoptions.SelectOptionsHandler) {
	selectOptionsRouter := r.Group("/select-options")
	selectOptionsRouter.GET("", handler.GetSelectOptions)
}
