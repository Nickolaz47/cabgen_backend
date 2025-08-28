package public_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testutils.SetupTestContext()

	db := testutils.SetupTestRepos()

	mockCountry := testmodels.NewCountry("", "", "", "")
	db.Create(&mockCountry)

	t.Run("Success", func(t *testing.T) {
		mockRegisterUser := testmodels.NewRegisterUser("", "")
		body := testutils.ToJSON(mockRegisterUser)
		expected := map[string]any{
			"message": "User created successfully. Please wait for activation by an administrator.",
			"data": map[string]any{
				"name":         mockRegisterUser.Name,
				"username":     mockRegisterUser.Username,
				"email":        mockRegisterUser.Email,
				"country_code": mockRegisterUser.CountryCode,
				"country":      "Brazil",
				"user_role":    "Collaborator",
			},
		}

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)
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
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/auth/register", tt.Body,
				nil, nil,
			)
			public.Register(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Email already exists", func(t *testing.T) {
		mockRegisterUser := testmodels.NewRegisterUser("nick", "")
		body := testutils.ToJSON(mockRegisterUser)
		expected := `{"error": "Email is already in use."}`

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)
		public.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Username already exists", func(t *testing.T) {
		mockRegisterUser := testmodels.NewRegisterUser("", "nick@mail.com")
		body := testutils.ToJSON(mockRegisterUser)
		expected := `{"error": "Username already exists."}`

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)
		public.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
