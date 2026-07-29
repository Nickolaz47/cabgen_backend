package data

import (
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

var baseValidRequestEmailUpdateBody = map[string]any{
	"new_email":         "new@example.com",
	"confirm_new_email": "new@example.com",
}

var RequestEmailUpdateTests = []Body{
	{"New email required", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidRequestEmailUpdateBody)
		b["new_email"] = ""
		return b
	}()), `{"error":"New email is required."}`},
	{"New email invalid", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidRequestEmailUpdateBody)
		b["new_email"] = "invalid-email"
		return b
	}()), `{"error":"Invalid email format."}`},
	{"Confirm new email required", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidRequestEmailUpdateBody)
		b["confirm_new_email"] = ""
		return b
	}()), `{"error":"Email confirmation is required."}`},
	{"Emails do not match", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidRequestEmailUpdateBody)
		b["confirm_new_email"] = "other@example.com"
		return b
	}()), `{"error":"Emails must match."}`},
}

var baseValidConfirmEmailUpdateBody = map[string]any{
	"token": "valid-token",
}

var ConfirmEmailUpdateTests = []Body{
	{"Token required", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidConfirmEmailUpdateBody)
		b["token"] = ""
		return b
	}()), `{"error":"Token is required."}`},
}
