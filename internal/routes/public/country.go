package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/country"
	"github.com/gin-gonic/gin"
)

func SetupCountryRoutes(r *gin.RouterGroup) {
	countryRouter := r.Group("/country")

	countryRouter.GET("", country.GetCountries)
	countryRouter.GET("/:code", country.GetCountryByID)
}
