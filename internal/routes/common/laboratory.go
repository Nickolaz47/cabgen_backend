package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/laboratory"
	"github.com/gin-gonic/gin"
)

func SetupLaboratoryRoutes(r *gin.RouterGroup, handler *laboratory.LaboratoryHandler) {
	laboratoryRouter := r.Group("/laboratory")
	laboratoryRouter.GET("", handler.GetActiveLaboratories)
}
