package public_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testutils.SetupTestContext()

	db := testutils.SetupTestRepos()

	db.Create(&models.Country{
		Code: "BRA",
		Pt:   "Brasil",
		En:   "Brazil",
		Es:   "Brazil",
	})
	t.Run("Success", func(t *testing.T) {
		body := `{"name":"Nicolas","username":"nmfaraujo","email":"nicolas@mail.com","confirm_email":"nicolas@mail.com","password":"12345678","confirm_password":"12345678","country_code":"BRA"}`
		expected := map[string]any{
			"message": "User created successfully. Please wait for activation by an administrator.",
			"data": map[string]any{
				"name":         "Nicolas",
				"username":     "nmfaraujo",
				"email":        "nicolas@mail.com",
				"country_code": "BRA",
				"country":      "Brazil",
				"user_role":    "Collaborator",
			},
		}

		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", body)
		public.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var got map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)

		if data, ok := got["data"].(map[string]any); ok {
			delete(data, "id")
		}

		assert.Equal(t, expected, got)
	})

	for _, tt := range data.RegisterTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", tt.Body)
			public.Register(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Email already exists", func(t *testing.T) {
		rawBody := map[string]string{
			"name":             "Nicolas",
			"username":         "nick",
			"email":            "nicolas@mail.com",
			"confirm_email":    "nicolas@mail.com",
			"password":         "12345678",
			"confirm_password": "12345678",
			"country_code":     "BRA",
			"user_role":        "Collaborator",
		}
		body := testutils.ToJSON(rawBody)
		expected := `{"error": "Email is already in use."}`

		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", body)
		public.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Username already exists", func(t *testing.T) {
		rawBody := map[string]string{
			"name":             "Nicolas",
			"username":         "nmfaraujo",
			"email":            "nick@mail.com",
			"confirm_email":    "nick@mail.com",
			"password":         "12345678",
			"confirm_password": "12345678",
			"country_code":     "BRA",
			"user_role":        "Collaborator",
		}
		body := testutils.ToJSON(rawBody)
		expected := `{"error": "Username already exists."}`

		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", body)
		public.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
func TestLogin(t *testing.T)   {}
func TestLogout(t *testing.T)  {}
func TestRefresh(t *testing.T) {}
