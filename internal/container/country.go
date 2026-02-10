package container

import (
	adminCountry "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/country"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/country"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildCountryService(db *gorm.DB, logger *zap.Logger) services.CountryService {
	countryRepo := repositories.NewCountryRepo(db)
	countryService := services.NewCountryService(countryRepo, logger)

	return countryService
}

func BuildPublicCountryHandler(svc services.CountryService) *country.PublicCountryHandler {
	return country.NewPublicCountryHandler(svc)
}

func BuildAdminCountryHandler(svc services.CountryService) *adminCountry.AdminCountryHandler {
	return adminCountry.NewAdminCountryHandler(svc)
}
