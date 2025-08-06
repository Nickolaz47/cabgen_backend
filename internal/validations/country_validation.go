package validations

import (
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

func ValidateCountryCode(countryCode string) (*models.Country, bool) {
	var country models.Country
	err := db.DB.First(&country, "code = ?", countryCode).Error
	if err != nil {
		return nil, false
	}
	return &country, true
}
