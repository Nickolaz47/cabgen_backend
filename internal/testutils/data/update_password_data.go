package data

import (
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

var baseValidUpdatePasswordBody = map[string]any{
	"current_password": "oldpass",
	"new_password":     "newpass123",
	"confirm_password": "newpass123",
}

var UpdatePasswordTests = []Body{
	{"Current password required", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidUpdatePasswordBody)
		b["current_password"] = ""
		return b
	}()), `{"error":"Current password is required."}`},
	{"New password required", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidUpdatePasswordBody)
		b["new_password"] = ""
		return b
	}()), `{"error":"New password is required."}`},
	{"New password too short", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidUpdatePasswordBody)
		b["new_password"] = "short"
		b["confirm_password"] = "short"
		return b
	}()), `{"error":"New password must be at least 8 characters."}`},
	{"New password too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidUpdatePasswordBody)
		b["new_password"] = strings.Repeat("a", 33)
		b["confirm_password"] = strings.Repeat("a", 33)
		return b
	}()), `{"error":"New password must be at most 32 characters."}`},
	{"Confirm password required", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidUpdatePasswordBody)
		b["confirm_password"] = ""
		return b
	}()), `{"error":"Password confirmation is required."}`},
	{"Passwords do not match", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidUpdatePasswordBody)
		b["confirm_password"] = "different123"
		return b
	}()), `{"error":"Passwords must match."}`},
}
