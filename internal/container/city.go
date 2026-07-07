package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/city"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func BuildCityHandler() *city.CityHandler {
	service := services.NewCityService()
	return city.NewCityHandler(service)
}
