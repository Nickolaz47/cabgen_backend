package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/gin-gonic/gin"
)

func SetupAdminLaboratoryRoutes(r *gin.RouterGroup, handler *laboratory.AdminLaboratoryHandler) {
	laboratoryRouter := r.Group("/laboratories")

	laboratoryRouter.GET("", handler.GetAllLaboratories)
	laboratoryRouter.GET("/:laboratoryId", handler.GetLaboratoryByID)
	laboratoryRouter.GET("/search", handler.GetLaboratoriesByNameOrAbbreviation)
	laboratoryRouter.POST("", handler.CreateLaboratory)
	laboratoryRouter.PUT("/:laboratoryId", handler.UpdateLaboratory)
	laboratoryRouter.DELETE("/:laboratoryId", handler.DeleteLaboratory)
}
