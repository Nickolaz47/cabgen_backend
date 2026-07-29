package data

import (
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

var baseValidLoginBody = map[string]any{
	"username": "nick",
	"password": "12345678",
}

var LoginSuccess = Body{
	Name:     "Success",
	Body:     testutils.ToJSON(baseValidLoginBody),
	Expected: `{"message": "Login successful."}`,
}

var LoginBadRequestTests = []Body{
	{"Username required", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseValidLoginBody); b["username"] = ""; return b }()), `{"error":"Username is required."}`},
	{"Password required", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseValidLoginBody); b["password"] = ""; return b }()), `{"error":"Password is required."}`},
}

var baseValidForgotBody = map[string]any{
	"email": "test@mail.com",
}

var ForgotPasswordBadRequestTests = []Body{
	{
		Name: "Email required",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseValidForgotBody)
			b["email"] = ""
			return b
		}()),
		Expected: `{"error":"Email is required."}`,
	},
	{
		Name: "Invalid Email Format",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseValidForgotBody)
			b["email"] = "not-an-email"
			return b
		}()),
		Expected: `{"error":"Invalid email format."}`,
	},
}

var baseValidResetBody = map[string]any{
	"token":            "valid-uuid-token-string",
	"new_password":     "newpassword123",
	"confirm_password": "newpassword123",
}

var ResetPasswordBadRequestTests = []Body{
	{
		Name: "Token required",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseValidResetBody)
			b["token"] = ""
			return b
		}()),
		Expected: `{"error":"Token is required."}`,
	},
	{
		Name: "New Password required",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseValidResetBody)
			b["new_password"] = ""
			return b
		}()),
		Expected: `{"error":"New password is required."}`,
	},
	{
		Name: "Confirm password required",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseValidResetBody)
			b["confirm_password"] = ""
			return b
		}()),
		Expected: `{"error":"Password confirmation is required."}`,
	},
	{
		Name: "Passwords do not match",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseValidResetBody)
			b["confirm_password"] = "different_password"
			return b
		}()),
		Expected: `{"error":"Passwords must match."}`,
	},
	{
		Name: "Password too short",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseValidResetBody)
			b["new_password"] = "short"
			b["confirm_password"] = "short"
			return b
		}()),
		Expected: `{"error":"New password must be at least 8 characters."}`,
	},
	{
		Name: "Password too long",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseValidResetBody)
			b["new_password"] = strings.Repeat("a", 33)
			b["confirm_password"] = strings.Repeat("a", 33)
			return b
		}()),
		Expected: `{"error":"New password must be at most 32 characters."}`,
	},
}
