package data

import (
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

var baseValidAdminUpdateBody = map[string]any{
	"name":         "Nicolas",
	"username":     "nmfaraujo",
	"country_code": "BRA",
	"interest":     "Bacterial resistance",
	"role":         "Researcher",
	"institution":  "NCBI",
}

var AdminUpdateUserTests = []Body{
	// Name
	{"Name too short", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseValidAdminUpdateBody); b["name"] = "Ni"; return b }()), `{"error":"Name must be at least 3 characters long."}`},
	{"Name too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminUpdateBody)
		b["name"] = strings.Repeat("nicolas", 15)
		return b
	}()), `{"error":"Name must be at most 100 characters long."}`},
	// Username
	{"Username too short", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminUpdateBody)
		b["username"] = "ni"
		return b
	}()), `{"error":"Username must be at least 4 characters long."}`},
	{"Username too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminUpdateBody)
		b["username"] = strings.Repeat("nick", 26)
		return b
	}()), `{"error":"Username must be at most 100 characters long."}`},
	// Email
	{"Email invalid", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminCreateUserBody)
		b["email"] = "invalid-email"
		return b
	}()), `{"error":"Invalid email format."}`},
	// Password
	{"Password too short", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminCreateUserBody)
		b["password"] = "123"
		return b
	}()), `{"error":"Password must be at least 8 characters long."}`},
	{"Password too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminCreateUserBody)
		b["password"] = strings.Repeat("1234", 10)
		return b
	}()), `{"error":"Password must be at most 32 characters long."}`},
	// User Role
	{"Invalid user role", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminCreateUserBody)
		b["user_role"] = "Tester"
		return b
	}()), `{"error":"The user role is invalid. Choose between Collaborator or Admin."}`},
	// Optional fields max
	{"Interest too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminUpdateBody)
		b["interest"] = string(make([]byte, 256))
		return b
	}()), `{"error":"Interest must be at most 255 characters long."}`},
	{"Role too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminUpdateBody)
		b["role"] = string(make([]byte, 256))
		return b
	}()), `{"error":"Role must be at most 255 characters long."}`},
	{"Institution too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidAdminUpdateBody)
		b["institution"] = string(make([]byte, 256))
		return b
	}()), `{"error":"Institution must be at most 255 characters long."}`},
}

var AdminCountryNotFoundTest = Body{
	"Country code invalid",
	testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidUpdateBody)
		b["country_code"] = "XXX"
		return b
	}()),
	`{"error":"No country was found with this code."}`,
}
