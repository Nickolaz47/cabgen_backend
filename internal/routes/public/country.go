package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/country"
	"github.com/gin-gonic/gin"
)

func SetupCountryRoutes(r *gin.RouterGroup, handler *country.PublicCountryHandler) {
	countryRouter := r.Group("/countries")

	countryRouter.GET("", handler.GetCountries)
}
