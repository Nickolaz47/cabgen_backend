package data

import "github.com/CABGenOrg/cabgen_backend/internal/testutils"

var baseLaboratoryCreateBody = map[string]any{
	"name":         "Laboratório Central do Rio de Janeiro",
	"abbreviation": "LACEN/RJ",
	"is_active":    true,
}

var CreateLaboratoryTests = []Body{
	{"Name Required", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseLaboratoryCreateBody); b["name"] = ""; return b }()), `{"error":"Name is required."}`},
	{"Name too short", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseLaboratoryCreateBody); b["name"] = "La"; return b }()), `{"error":"Name must be at least 3 characters long."}`},
	{"Abbreviation Required", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseLaboratoryCreateBody)
		b["abbreviation"] = ""
		return b
	}()), `{"error":"Abbreviation is required."}`},
	{"Abbreviation too short", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseLaboratoryCreateBody)
		b["abbreviation"] = "L"
		return b
	}()), `{"error":"Abbreviation must be at least 2 characters long."}`},
}

var UpdateLaboratoryTests = []Body{
	{"Name too short", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseLaboratoryCreateBody); b["name"] = "La"; return b }()), `{"error":"Name must be at least 3 characters long."}`},
	{"Abbreviation too short", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseLaboratoryCreateBody)
		b["abbreviation"] = "L"
		return b
	}()), `{"error":"Abbreviation must be at least 2 characters long."}`},
}
