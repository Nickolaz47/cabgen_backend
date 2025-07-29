package country

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/country"
	"github.com/gin-gonic/gin"
)

func CountryRoutes(r *gin.Engine) {
	r.GET("/country", country.GetCountries)
	r.GET("/country/:code", country.GetCountryByID)
}
