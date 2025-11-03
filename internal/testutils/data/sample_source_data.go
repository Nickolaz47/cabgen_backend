package data

import "github.com/CABGenOrg/cabgen_backend/internal/testutils"

var baseSampleSourceCreateBody = map[string]any{
	"names":     map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
	"groups":    map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
	"is_active": true,
}

var CreateSampleSourceTests = []Body{
	{"Missing names", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseSampleSourceCreateBody); b["names"] = nil; return b }()), `{"error":"The names parameter with translations for pt, en, and es is required."}`},
	{"Invalid names", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceCreateBody)
		b["names"] = map[string]string{"pt": "Plasma"}
		return b
	}()), `{"error":"The names parameter must contain at least 3 keys (pt, en, es)."}`},
	{"Missing names key", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceCreateBody)
		b["names"] = map[string]string{"pt": "Plasma", "fr": "Plasma", "es": "Plasma"}
		return b
	}()), `{"error":"Sample source missing translation for en."}`},
	{"Empty name translation", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceCreateBody)
		b["names"] = map[string]string{"pt": "Plasma", "en": "", "es": "Plasma"}
		return b
	}()), `{"error":"Sample source has empty translation for en."}`},
	{"Missing groups", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseSampleSourceCreateBody); b["groups"] = nil; return b }()), `{"error":"The groups parameter with translations for pt, en, and es is required."}`},
	{"Invalid groups", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceCreateBody)
		b["groups"] = map[string]string{"pt": "Sangue"}
		return b
	}()), `{"error":"The groups parameter must contain at least 3 keys (pt, en, es)."}`},
	{"Missing groups key", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceCreateBody)
		b["groups"] = map[string]string{"pt": "Sangue", "fr": "Sang", "es": "Sangre"}
		return b
	}()), `{"error":"Sample source missing translation for en."}`},
	{"Empty group translation", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceCreateBody)
		b["groups"] = map[string]string{"pt": "Sangue", "en": "", "es": "Sangre"}
		return b
	}()), `{"error":"Sample source has empty translation for en."}`},
}

var baseSampleSourceUpdateBody = map[string]any{
	"names":     map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
	"groups":    map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
	"is_active": true,
}

var UpdateSampleSourceTests = []Body{
	{"Invalid names", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceUpdateBody)
		b["names"] = map[string]string{"pt": "Aspirado"}
		return b
	}()), `{"error":"The names parameter must contain at least 3 keys (pt, en, es)."}`},
	{"Missing names key", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceUpdateBody)
		b["names"] = map[string]string{"pt": "Aspirado", "fr": "Aspiré", "es": "Aspirado"}
		return b
	}()), `{"error":"Sample source has empty translation for en."}`},
	{"Empty name translation", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceUpdateBody)
		b["names"] = map[string]string{"pt": "Aspirado", "en": "", "es": "Aspirado"}
		return b
	}()), `{"error":"Empty en translation."}`},
	{"Invalid groups", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceUpdateBody)
		b["groups"] = map[string]string{"pt": "Trato respiratório"}
		return b
	}()), `{"error":"The names parameter must contain at least 3 keys (pt, en, es)."}`},
	{"Missing groups key", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceUpdateBody)
		b["groups"] = map[string]string{"pt": "Trato respiratório", "fr": "Voies respiratoires", "es": "Vías respiratorias"}
		return b
	}()), `{"error":"Missing en translation."}`},
	{"Empty group translation", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleSourceUpdateBody)
		b["groups"] = map[string]string{"pt": "Trato respiratório", "en": "", "es": "Vías respiratorias"}
		return b
	}()), `{"error":"Empty en translation."}`},
}
