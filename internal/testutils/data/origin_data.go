package data

import "github.com/CABGenOrg/cabgen_backend/internal/testutils"

var baseOriginCreateBody = map[string]any{
	"names":     map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
	"is_active": true,
}

var CreateOriginTests = []Body{
	{"Missing names", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseOriginCreateBody); b["names"] = nil; return b }()), `{"error":"The names parameter with translations for pt, en, and es is required."}`},
	{"Invalid names", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseOriginCreateBody)
		b["names"] = map[string]string{"pt": "Humano"}
		return b
	}()), `{"error":"The names parameter must contain at least 3 keys (pt, en, es)."}`},
	{"Missing names key", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseOriginCreateBody)
		b["names"] = map[string]string{"pt": "Alimentar", "fr": "Nourrir", "es": "Alimentaria"}
		return b
	}()), `{"error":"Missing en translation."}`},
	{"Empty translation", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseOriginCreateBody)
		b["names"] = map[string]string{"pt": "Alimentar", "en": "", "es": "Alimentaria"}
		return b
	}()), `{"error":"Empty en translation."}`},
}

var baseOriginUpdateBody = map[string]any{
	"names":     map[string]string{"pt": "Ambiental", "en": "Environmental", "es": "Ambiental"},
	"is_active": true,
}

var UpdateOriginTests = []Body{
	{"Invalid names", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseOriginUpdateBody)
		b["names"] = map[string]string{"pt": "Humano"}
		return b
	}()), `{"error":"The names parameter must contain at least 3 keys (pt, en, es)."}`},
	{"Missing names key", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseOriginUpdateBody)
		b["names"] = map[string]string{"pt": "Alimentar", "fr": "Nourrir", "es": "Alimentaria"}
		return b
	}()), `{"error":"Missing en translation."}`},
	{"Empty translation", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseOriginUpdateBody)
		b["names"] = map[string]string{"pt": "Alimentar", "en": "", "es": "Alimentaria"}
		return b
	}()), `{"error":"Empty en translation."}`},
}
