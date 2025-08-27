package models

import "github.com/CABGenOrg/cabgen_backend/internal/models"

func NewCountry(code, pt, en, es string) models.Country {
	if code == "" {
		code = "BRA"
	}

	if pt == "" {
		pt = "Brasil"
	}

	if en == "" {
		en = "Brazil"
	}

	if es == "" {
		es = "Brazil"
	}

	return models.Country{
		Code: code,
		Pt:   pt,
		En:   en,
		Es:   es,
	}
}
