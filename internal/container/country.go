package container

import (
	adminCountry "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/country"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/country"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildCountryService(db *gorm.DB) services.CountryService {
	countryRepo := repository.NewCountryRepo(db)
	countryService := services.NewCountryService(countryRepo)

	return countryService
}

func BuildPublicCountryHandler(svc services.CountryService) *country.PublicCountryHandler {
	return country.NewPublicCountryHandler(svc)
}

func BuildAdminCountryHandler(svc services.CountryService) *adminCountry.AdminCountryHandler {
	return adminCountry.NewAdminCountryHandler(svc)
}
