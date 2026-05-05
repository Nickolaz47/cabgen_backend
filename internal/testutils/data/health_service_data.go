package data

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

var baseHealthServiceBody = map[string]any{
	"name":          "Hospital de Clínicas",
	"type":          models.Private,
	"country_code":  "BRA",
	"city":          "Rio de Janeiro",
	"contactant":    "Dr. Roberto",
	"contact_email": "contato@hospital.com",
	"contact_phone": "+5521999999999",
	"is_active":     true,
}

var CreateHealthServiceTests = []Body{
	{
		"Name Required",
		testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseHealthServiceBody); b["name"] = ""; return b }()),
		`{"error":"Name is required."}`,
	},
	{
		"Name too short",
		testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseHealthServiceBody); b["name"] = "Ho"; return b }()),
		`{"error":"Name must be at least 3 characters long."}`,
	},
	{
		"Type Required",
		testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseHealthServiceBody); b["type"] = ""; return b }()),
		`{"error":"The health service type is required."}`,
	},
	{
		"Country Code Required",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["country_code"] = ""
			return b
		}()),
		`{"error":"Country code is required."}`,
	},
	{
		"Country Code invalid length",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["country_code"] = "BR"
			return b
		}()),
		`{"error":"The country code must be 3 characters long."}`,
	},
	{
		"Contactant name too short",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["contactant"] = "Dr"
			return b
		}()),
		`{"error":"Contact name must be at least 3 characters long."}`,
	},
	{
		"Invalid Email format",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["contact_email"] = "email-invalido"
			return b
		}()),
		`{"error":"The contact email format is invalid."}`,
	},
	{
		"Invalid Phone format (E164)",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["contact_phone"] = "021999999999"
			return b
		}()),
		`{"error":"Contact phone number is in an invalid format."}`,
	},
}

var UpdateHealthServiceTests = []Body{
	{
		"Name too short",
		testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseHealthServiceBody); b["name"] = "Ho"; return b }()),
		`{"error":"Name must be at least 3 characters long."}`,
	},
	{
		"Country Code invalid length",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["country_code"] = "BR"
			return b
		}()),
		`{"error":"The country code must be 3 characters long."}`,
	},
	{
		"Contactant name too short",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["contactant"] = "Dr"
			return b
		}()),
		`{"error":"Contact name must be at least 3 characters long."}`,
	},
	{
		"Invalid Email format",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["contact_email"] = "email-invalido"
			return b
		}()),
		`{"error":"The contact email format is invalid."}`,
	},
	{
		"Invalid Phone format (E164)",
		testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseHealthServiceBody)
			b["contact_phone"] = "021999999999"
			return b
		}()),
		`{"error":"Contact phone number is in an invalid format."}`,
	},
}
