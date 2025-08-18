package data

import (
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

var baseValidRegisterBody = map[string]string{
	"name":             "Eddie",
	"username":         "eddy",
	"email":            "eddy@mail.com",
	"confirm_email":    "eddy@mail.com",
	"password":         "12345678",
	"confirm_password": "12345678",
	"country_code":     "BRA",
	"interest":         "Bacterial resistance",
	"role":             "Researcher",
	"institution":      "NCBI",
}

var RegisterTests = []Body{
	// Name
	{"Name required", testutils.ToJSON(func() map[string]string { b := testutils.CopyMap(baseValidRegisterBody); b["name"] = ""; return b }()), `{"error":"Name is required."}`},
	{"Name too short", testutils.ToJSON(func() map[string]string { b := testutils.CopyMap(baseValidRegisterBody); b["name"] = "Ed"; return b }()), `{"error":"Name must be at least 3 characters long."}`},
	{"Name too long", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["name"] = strings.Repeat("eddie", 21)
		return b
	}()), `{"error":"Name must be at most 100 characters long."}`},

	// Username
	{"Username required", testutils.ToJSON(func() map[string]string { b := testutils.CopyMap(baseValidRegisterBody); b["username"] = ""; return b }()), `{"error":"Username is required."}`},
	{"Username too short", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["username"] = "ed"
		return b
	}()), `{"error":"Username must be at least 4 characters long."}`},
	{"Username too long", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["username"] = strings.Repeat("eddy", 26)
		return b
	}()), `{"error":"Username must be at most 100 characters long."}`},

	// Email
	{"Email required", testutils.ToJSON(func() map[string]string { b := testutils.CopyMap(baseValidRegisterBody); b["email"] = ""; return b }()), `{"error":"Email is required."}`},
	{"Email invalid", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["email"] = "invalid-email"
		return b
	}()), `{"error":"Invalid email format."}`},
	{"Emails not match", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["confirm_email"] = "other@mail.com"
		return b
	}()), `{"error":"Emails must match."}`},

	// Password
	{"Password required", testutils.ToJSON(func() map[string]string { b := testutils.CopyMap(baseValidRegisterBody); b["password"] = ""; return b }()), `{"error":"Password is required."}`},
	{"Password too short", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["password"] = "123"
		return b
	}()), `{"error":"Password must be at least 8 characters long."}`},
	{"Password too long", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["password"] = strings.Repeat("1234", 10)
		return b
	}()), `{"error":"Password must be at most 32 characters long."}`},
	{"Passwords not match", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["confirm_password"] = "87654321"
		return b
	}()), `{"error":"Passwords must match."}`},

	// Country code
	{"Country code required", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["country_code"] = ""
		return b
	}()), `{"error":"Country code is required."}`},
	{"Country code invalid", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["country_code"] = "XXX"
		return b
	}()), `{"error":"No country was found with this code."}`},

	// Optional fields max
	{"Interest too long", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["interest"] = string(make([]byte, 256))
		return b
	}()), `{"error":"Interest must be at most 255 characters long."}`},
	{"Role too long", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["role"] = string(make([]byte, 256))
		return b
	}()), `{"error":"Role must be at most 255 characters long."}`},
	{"Institution too long", testutils.ToJSON(func() map[string]string {
		b := testutils.CopyMap(baseValidRegisterBody)
		b["institution"] = string(make([]byte, 256))
		return b
	}()), `{"error":"Institution must be at most 255 characters long."}`},
}
