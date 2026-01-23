package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/country"
	"github.com/gin-gonic/gin"
)

func SetupAdminCountryRoutes(r *gin.RouterGroup, handler *country.AdminCountryHandler) {
	originRouter := r.Group("/countries")

	originRouter.GET("", handler.GetCountries)
	originRouter.GET("/:code", handler.GetCountryByCode)
	originRouter.GET("/search", handler.GetCountriesByName)
	originRouter.POST("", handler.CreateCountry)
	originRouter.PUT("/:code", handler.UpdateCountry)
	originRouter.DELETE("/:code", handler.DeleteCountry)
}
