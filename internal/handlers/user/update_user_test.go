package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	mockUserToken := testmodels.NewUserToken(
		mockLoginUser.ID, mockLoginUser.Username, mockLoginUser.UserRole,
	)

	gc, gw := testutils.SetupGinContext(
		http.MethodPost,
		"/api/auth/login",
		testutils.ToJSON(models.LoginInput{
			Username: mockLoginUser.Username,
			Password: "12345678"},
		), nil, nil,
	)

	public.Login(gc)
	var accessCookie string

	cookies := gw.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == auth.Access {
			accessCookie = cookie.Value
		}
	}

	t.Run("Missing user in context", func(t *testing.T) {
		mockUpdateUserInput := testmodels.NewUpdateUserInput()
		body := testutils.ToJSON(mockUpdateUserInput)
		req := httptest.NewRequest(http.MethodPut, "/api/user/me", bytes.NewBufferString(body))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		user.UpdateUser(c)

		expected := `{"error": "Unauthorized. Please log in to continue."}`

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateUserTests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/user/me", bytes.NewBufferString(tt.Body))
			req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			c.Set("user", &mockUserToken)

			user.UpdateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run(data.CountryNotFoundTest.Name, func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/user/me", bytes.NewBufferString(data.CountryNotFoundTest.Body))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		c.Set("user", &mockUserToken)

		user.UpdateUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, data.CountryNotFoundTest.Expected, w.Body.String())
	})

	for _, tt := range data.UpdateUserConflictTests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/user/me", bytes.NewBufferString(tt.Body))
			req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			c.Set("user", &mockUserToken)

			user.UpdateUser(c)

			assert.Equal(t, http.StatusConflict, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("User not found", func(t *testing.T) {
		mockUpdateUserInput := testmodels.NewUpdateUserInput()
		body := testutils.ToJSON(mockUpdateUserInput)
		req := httptest.NewRequest(http.MethodPut, "/api/user/me", bytes.NewBufferString(body))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		c.Set("user", &models.UserToken{
			ID:       uuid.UUID{},
			Username: mockLoginUser.Username,
			UserRole: mockLoginUser.UserRole,
		})

		user.UpdateUser(c)

		expected := `{"error": "User not found."}`

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		mockUpdateUserInput := testmodels.NewUpdateUserInput()
		body := testutils.ToJSON(mockUpdateUserInput)
		req := httptest.NewRequest(http.MethodPut, "/api/user/me", bytes.NewBufferString(body))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		c.Set("user", &mockUserToken)

		user.UpdateUser(c)

		expected := map[string]any{
			"data": map[string]any{
				"name":         *mockUpdateUserInput.Name,
				"username":     *mockUpdateUserInput.Username,
				"email":        mockLoginUser.Email,
				"country_code": *mockUpdateUserInput.CountryCode,
				"country":      "Brazil",
				"user_role":    "",
				"role":         *mockUpdateUserInput.Role,
				"interest":     *mockUpdateUserInput.Interest,
				"institution":  *mockUpdateUserInput.Institution,
			},
		}

		assert.Equal(t, http.StatusOK, w.Code)

		var got map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)

		if data, ok := got["data"].(map[string]any); ok {
			delete(data, "id")
		}

		assert.Equal(t, expected, got)
	})
}
