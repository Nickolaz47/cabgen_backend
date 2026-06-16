package data

import (
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

var baseTicketCreateBody = map[string]any{
	"name":        "Jão",
	"email":       "jao@mail.com",
	"institution": "Fiocruz",
	"subject":     "Wrong password",
	"message":     "Cannot access my account.",
}

var CreateTicketValidationTests = []Body{
	{
		Name: "Missing name",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["name"] = ""
			return b
		}()),
		Expected: `{"error":"Name is required."}`,
	},
	{
		Name: "Name too short",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["name"] = "A"
			return b
		}()),
		Expected: `{"error":"Name must be at least 2 characters long."}`,
	},
	{
		Name: "Name too long",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["name"] = strings.Repeat("A", 101)
			return b
		}()),
		Expected: `{"error":"Name must be at most 100 characters long."}`,
	},
	{
		Name: "Missing email",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["email"] = ""
			return b
		}()),
		Expected: `{"error":"Email is required."}`,
	},
	{
		Name: "Invalid email format",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["email"] = "invalid-email"
			return b
		}()),
		Expected: `{"error":"Invalid email format."}`,
	},
	{
		Name: "Missing institution",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["institution"] = ""
			return b
		}()),
		Expected: `{"error":"Institution is required."}`,
	},
	{
		Name: "Institution too long",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["institution"] = strings.Repeat("A", 151)
			return b
		}()),
		Expected: `{"error":"Institution must be at most 150 characters long."}`,
	},
	{
		Name: "Missing subject",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["subject"] = ""
			return b
		}()),
		Expected: `{"error":"Contact subject is required."}`,
	},
	{
		Name: "Subject too short",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["subject"] = "Abcd"
			return b
		}()),
		Expected: `{"error":"Contact subject must be at least 5 characters long."}`,
	},
	{
		Name: "Subject too long",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["subject"] = strings.Repeat("A", 151)
			return b
		}()),
		Expected: `{"error":"Contact subject must be at most 150 characters long."}`,
	},
	{
		Name: "Missing message",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["message"] = ""
			return b
		}()),
		Expected: `{"error":"Message is required."}`,
	},
	{
		Name: "Message too short",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["message"] = "Abcdefghi"
			return b
		}()),
		Expected: `{"error":"Message must be at least 10 characters long."}`,
	},
	{
		Name: "Message too long",
		Body: testutils.ToJSON(func() map[string]any {
			b := testutils.CopyMap(baseTicketCreateBody)
			b["message"] = strings.Repeat("A", 2001)
			return b
		}()),
		Expected: `{"error":"Message must be at most 2000 characters long."}`,
	},
}
