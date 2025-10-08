package user_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockAdminLoginUser := testmodels.NewAdminLoginUser()
	db.Create(&mockAdminLoginUser)

	mockAdminUserToken := testmodels.NewUserToken(
		mockAdminLoginUser.ID,
		mockAdminLoginUser.Username,
		mockAdminLoginUser.UserRole,
	)

	gc, _ := testutils.SetupGinContext(
		http.MethodPost,
		"/api/auth/login",
		testutils.ToJSON(models.LoginInput{
			Username: mockAdminLoginUser.Username,
			Password: "12345678"},
		), nil, nil,
	)

	public.Login(gc)

	t.Run("Missing user in context", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/user", "",
			nil, nil,
		)

		user.CreateUser(c)

		expected := `{"error": "Unauthorized. Please log in to continue."}`

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.AdminCreateUserTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/user", tt.Body,
				nil, nil,
			)

			c.Set("user", &mockAdminUserToken)

			user.CreateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Email already exists", func(t *testing.T) {
		body := testutils.ToJSON(
			testmodels.NewAdminCreateUserInput("cadmin@mail.com", ""),
		)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/user", body,
			nil, nil,
		)

		c.Set("user", &mockAdminUserToken)

		user.CreateUser(c)

		expected := `{"error": "Email is already in use."}`

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Username already exists", func(t *testing.T) {
		body := testutils.ToJSON(
			testmodels.NewAdminCreateUserInput("", "admin"),
		)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/user", body,
			nil, nil,
		)

		c.Set("user", &mockAdminUserToken)

		user.CreateUser(c)

		expected := `{"error": "Username already exists."}`

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		mockAdminCreateInput := testmodels.NewAdminCreateUserInput("", "")
		body := testutils.ToJSON(mockAdminCreateInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/user", body,
			nil, nil,
		)

		c.Set("user", &mockAdminUserToken)
		user.CreateUser(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "User created successfully.",
			"data": map[string]any{
				"name":         mockAdminCreateInput.Name,
				"username":     mockAdminCreateInput.Username,
				"email":        mockAdminCreateInput.Email,
				"country_code": mockAdminCreateInput.CountryCode,
				"country":      "Brazil",
				"user_role":    "Collaborator",
			},
		})

		var got map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)

		if data, ok := got["data"].(map[string]any); ok {
			delete(data, "id")
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, testutils.ToJSON(got))
	})
}
