package models

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

func NewCountry(code string, names map[string]string) models.Country {
	if code == "" {
		code = "BRA"
	}

	if names == nil {
		names = map[string]string{
			"pt": "Brasil",
			"en": "Brazil",
			"es": "Brazil",
		}
	}

	return models.Country{
		Code:  code,
		Names: names,
	}
}
