package data

import "github.com/CABGenOrg/cabgen_backend/internal/testutils"
import "github.com/CABGenOrg/cabgen_backend/internal/models"

var baseMicroorganismCreateBody = map[string]any{
	"taxon":     models.Bacteria,
	"species":   "Escherichia coli",
	"variety":   map[string]string{"pt": "Padrão", "en": "Standard", "es": "Estándar"},
	"is_active": true,
}

var CreateMicroorganismTests = []Body{
	{"Missing species", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismCreateBody)
		b["species"] = nil
		return b
	}()), `{"error":"Species is required."}`},

	{"Invalid taxon", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismCreateBody)
		b["taxon"] = "alien_life_form"
		return b
	}()), `{"error":"The microorganism taxon is invalid. Choose between Bacteria, Virus, Protozoa, or Fungi."}`},

	{"Invalid variety keys", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismCreateBody)
		b["variety"] = map[string]string{"pt": "Padrão"}
		return b
	}()), `{"error":"The variety parameter must contain at least 3 keys (pt, en, es)."}`},

	{"Missing variety key", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismCreateBody)
		b["variety"] = map[string]string{"pt": "Padrão", "fr": "Standard", "es": "Estándar"}
		return b
	}()), `{"error":"Microorganism variety missing translation for en."}`},

	{"Empty translation in variety", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismCreateBody)
		b["variety"] = map[string]string{"pt": "Padrão", "en": "", "es": "Estándar"}
		return b
	}()), `{"error":"Microorganism variety has empty translation for en."}`},
}

var baseMicroorganismUpdateBody = map[string]any{
	"taxon":   models.Virus,
	"species": "Influenza A",
	"variety": map[string]string{"pt": "H1N1", "en": "H1N1", "es": "H1N1"},
}

var UpdateMicroorganismTests = []Body{
	{"Invalid taxon", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismUpdateBody)
		b["taxon"] = "invalid_taxon"
		return b
	}()), `{"error":"The microorganism taxon is invalid. Choose between Bacteria, Virus, Protozoa, or Fungi."}`},

	{"Invalid variety keys", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismUpdateBody)
		b["variety"] = map[string]string{"pt": "H1N1"}
		return b
	}()), `{"error":"The variety parameter must contain at least 3 keys (pt, en, es)."}`},

	{"Missing variety key", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismUpdateBody)
		b["variety"] = map[string]string{"pt": "H1N1", "fr": "H1N1", "es": "H1N1"}
		return b
	}()), `{"error":"Microorganism variety missing translation for en."}`},

	{"Empty translation in variety", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseMicroorganismUpdateBody)
		b["variety"] = map[string]string{"pt": "H1N1", "en": "", "es": "H1N1"}
		return b
	}()), `{"error":"Microorganism variety has empty translation for en."}`},
}
