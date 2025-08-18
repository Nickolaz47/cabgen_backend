package validations

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/country"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

func ValidateCountryCode(countryCode string) (*models.Country, bool) {
	country, err := country.CountryRepo.GetCountry(countryCode)
	if err != nil {
		return nil, false
	}

	return country, true
}
