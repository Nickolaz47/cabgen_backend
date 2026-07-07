package common

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/city"
	"github.com/gin-gonic/gin"
)

func SetupCityRoutes(r *gin.RouterGroup,
	handler *city.CityHandler) {
	selectOptionsRouter := r.Group("/cities")
	selectOptionsRouter.GET("", handler.GetCities)
}
