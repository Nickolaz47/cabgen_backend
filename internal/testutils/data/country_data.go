package data

import "github.com/CABGenOrg/cabgen_backend/internal/testutils"

var baseCountryCreateBody = map[string]any{
	"code":  "BRA",
	"names": map[string]string{"pt": "Brasil", "en": "Brazil", "es": "Brasil"},
}

var CreateCountryTests = []Body{
	{
		"Missing code",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryCreateBody)
			b["code"] = nil
			return b
		}()),
		`{"error":"Country code is required."}`,
	},
	{
		"Invalid code length",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryCreateBody)
			b["code"] = "BR"
			return b
		}()),
		`{"error":"The country code must be 3 characters long."}`,
	},
	{
		"Missing names",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryCreateBody)
			b["names"] = nil
			return b
		}()),
		`{"error":"The names parameter with translations for pt, en, and es is required."}`,
	},
	{
		"Invalid names",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryCreateBody)
			b["names"] = map[string]string{"pt": "Brasil"}
			return b
		}()),
		`{"error":"The names parameter must contain at least 3 keys (pt, en, es)."}`,
	},
	{
		"Missing names key",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryCreateBody)
			b["names"] = map[string]string{"pt": "Brasil", "fr": "Brésil", "es": "Brasil"}
			return b
		}()),
		`{"error":"Country missing translation for en."}`,
	},
	{
		"Empty translation",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryCreateBody)
			b["names"] = map[string]string{"pt": "Brasil", "en": "", "es": "Brasil"}
			return b
		}()),
		`{"error":"Country has empty translation for en."}`,
	},
}

var baseCountryUpdateBody = map[string]any{
	"code":  "BRZ",
	"names": map[string]string{"pt": "Brasil", "en": "Brazil", "es": "Brasil"},
}

var UpdateCountryTests = []Body{
	{
		"Invalid code length",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryUpdateBody)
			b["code"] = "BR"
			return b
		}()),
		`{"error":"The country code must be 3 characters long."}`,
	},
	{
		"Invalid names",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryUpdateBody)
			b["names"] = map[string]string{"pt": "Brasil"}
			return b
		}()),
		`{"error":"The names parameter must contain at least 3 keys (pt, en, es)."}`,
	},
	{
		"Missing names key",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryUpdateBody)
			b["names"] = map[string]string{"pt": "Brasil", "fr": "Brésil", "es": "Brasil"}
			return b
		}()),
		`{"error":"Country missing translation for en."}`,
	},
	{
		"Empty translation",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseCountryUpdateBody)
			b["names"] = map[string]string{"pt": "Brasil", "en": "", "es": "Brasil"}
			return b
		}()),
		`{"error":"Country has empty translation for en."}`,
	},
}
