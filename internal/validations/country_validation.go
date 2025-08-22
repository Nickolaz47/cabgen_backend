package validations

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
)

func ValidateCountryCode(countryCode string) (*models.Country, bool) {
	country, err := repository.CountryRepo.GetCountry(countryCode)
	if err != nil {
		return nil, false
	}

	return country, true
}
