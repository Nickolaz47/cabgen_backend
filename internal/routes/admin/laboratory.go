package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/gin-gonic/gin"
)

func SetupLaboratoryRoutes(r *gin.RouterGroup, handler *laboratory.AdminLaboratoryHandler) {
	laboratoryRouter := r.Group("/laboratory")

	laboratoryRouter.GET("", handler.GetAllLaboratories)
	laboratoryRouter.GET("/:laboratoryId", handler.GetLaboratoryByID)
	laboratoryRouter.GET("/search", handler.GetLaboratoryByNameOrAbbreviation)
	laboratoryRouter.POST("", handler.CreateLaboratory)
	laboratoryRouter.PUT("/:laboratoryId", handler.UpdateLaboratory)
	laboratoryRouter.DELETE("/:laboratoryId", handler.DeleteLaboratory)
}
