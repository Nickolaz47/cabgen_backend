package user_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	t.Run("User not found", func(t *testing.T) {
		body := testutils.ToJSON(testmodels.NewAdminUpdateUserInput())
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/user/nmfaraujo", body,
			nil, gin.Params{{Key: "username", Value: "nmfaraujo2"}},
		)

		user.UpdateUser(c)

		expected := `{"error": "User not found."}`

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.AdminUpdateUserTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/user/nmfaraujo", tt.Body,
				nil, nil,
			)

			user.UpdateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run(data.AdminCountryNotFoundTest.Name, func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/user/nick", data.AdminCountryNotFoundTest.Body,
			nil, gin.Params{{Key: "username", Value: "nick"}},
		)

		user.UpdateUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, data.AdminCountryNotFoundTest.Expected, w.Body.String())
	})

	for _, tt := range data.AdminUpdateUserConflictTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/user/nick", tt.Body,
				nil, gin.Params{{Key: "username", Value: "nick"}},
			)

			user.UpdateUser(c)

			assert.Equal(t, http.StatusConflict, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Success", func(t *testing.T) {
		mockAdminUpdateUserInput := testmodels.NewAdminUpdateUserInput()
		body := testutils.ToJSON(mockAdminUpdateUserInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/user/nick", body,
			nil, gin.Params{{Key: "username", Value: "nick"}},
		)

		user.UpdateUser(c)

		var got map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)

		if data, ok := got["data"].(map[string]any); ok {
			delete(data, "id")
		}

		result := testutils.ToJSON(got)

		expected := testutils.ToJSON(
			map[string]map[string]any{
				"data": {
					"name":         *mockAdminUpdateUserInput.Name,
					"username":     *mockAdminUpdateUserInput.Username,
					"email":        *mockAdminUpdateUserInput.Email,
					"country_code": *mockAdminUpdateUserInput.CountryCode,
					"country":      "Brazil",
					"user_role":    *mockAdminUpdateUserInput.UserRole,
					"role":         *mockAdminUpdateUserInput.Role,
					"interest":     *mockAdminUpdateUserInput.Interest,
					"institution":  *mockAdminUpdateUserInput.Institution,
				}},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, result)
	})
}
