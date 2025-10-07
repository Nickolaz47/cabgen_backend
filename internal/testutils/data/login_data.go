package data

import "github.com/CABGenOrg/cabgen_backend/internal/testutils"

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

var LoginUnauthorizedTests = []Body{
	{"Invalid credentials (username)", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseValidLoginBody); b["username"] = "nic"; return b }()), `{"error":"Invalid credentials."}`},
	{"Invalid credentials (password)", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseValidLoginBody)
		b["password"] = "1234567"
		return b
	}()), `{"error":"Invalid credentials."}`},
}

var LoginUserDeactivatedTest = Body{
	Name:     "User deactivated",
	Body:     testutils.ToJSON(baseValidLoginBody),
	Expected: `{"error": "Your account is not activated."}`,
}
