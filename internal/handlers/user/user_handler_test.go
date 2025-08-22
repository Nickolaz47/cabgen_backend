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
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetOwnUser(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	db.Create(&models.Country{
		Code: "BRA",
		Pt:   "Brasil",
		En:   "Brazil",
		Es:   "Brazil",
	})
	db.Create(&testmodels.MockLoginUser)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/auth/login",
			testutils.ToJSON(models.LoginInput{
				Username: testmodels.MockLoginUser.Username,
				Password: testmodels.MockRegisterUser.Password},
			),
		)

		public.Login(c)
		var accessCookie string

		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == auth.Access {
				accessCookie = cookie.Value
			}
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/me", bytes.NewBuffer(nil))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		c.Set("user", &models.UserToken{
			ID:       testmodels.MockLoginUser.ID,
			Username: testmodels.MockLoginUser.Username,
			UserRole: testmodels.MockLoginUser.UserRole,
		})

		user.GetOwnUser(c)

		expected := testutils.ToJSON(responses.APIResponse{Data: testmodels.MockLoginUser.ToResponse(c)})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Missing user in context", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/user/me",
			"",
		)

		user.GetOwnUser(c)

		expected := `{"error": "Unauthorized. Please log in to continue."}`

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("User not found", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/user/me",
			"",
		)

		mockUserToken := &models.UserToken{
			ID:       uuid.UUID{},
			Username: "nick",
			UserRole: "Collaborator",
		}
		c.Set("user", mockUserToken)

		user.GetOwnUser(c)

		expected := `{"error": "Unauthorized. Please log in to continue."}`

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}

func TestUpdateUser(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	db.Create(&models.Country{
		Code: "BRA",
		Pt:   "Brasil",
		En:   "Brazil",
		Es:   "Brazil",
	})
	db.Create(&testmodels.MockLoginUser)

	gc, gw := testutils.SetupGinContext(
		http.MethodPost,
		"/api/auth/login",
		testutils.ToJSON(models.LoginInput{
			Username: testmodels.MockLoginUser.Username,
			Password: testmodels.MockRegisterUser.Password},
		),
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
		body := testutils.ToJSON(testmodels.MockUpdateUser)
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

			c.Set("user", &models.UserToken{
				ID:       testmodels.MockLoginUser.ID,
				Username: testmodels.MockLoginUser.Username,
				UserRole: testmodels.MockLoginUser.UserRole,
			})

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

		c.Set("user", &models.UserToken{
			ID:       testmodels.MockLoginUser.ID,
			Username: testmodels.MockLoginUser.Username,
			UserRole: testmodels.MockLoginUser.UserRole,
		})

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

			c.Set("user", &models.UserToken{
				ID:       testmodels.MockLoginUser.ID,
				Username: testmodels.MockLoginUser.Username,
				UserRole: testmodels.MockLoginUser.UserRole,
			})

			user.UpdateUser(c)

			assert.Equal(t, http.StatusConflict, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("User not found", func(t *testing.T) {
		body := testutils.ToJSON(testmodels.MockUpdateUser)
		req := httptest.NewRequest(http.MethodPut, "/api/user/me", bytes.NewBufferString(body))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		c.Set("user", &models.UserToken{
			ID:       uuid.UUID{},
			Username: testmodels.MockLoginUser.Username,
			UserRole: testmodels.MockLoginUser.UserRole,
		})

		user.UpdateUser(c)

		expected := `{"error": "User not found."}`

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		body := testutils.ToJSON(testmodels.MockUpdateUser)
		req := httptest.NewRequest(http.MethodPut, "/api/user/me", bytes.NewBufferString(body))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		c.Set("user", &models.UserToken{
			ID:       testmodels.MockLoginUser.ID,
			Username: testmodels.MockLoginUser.Username,
			UserRole: testmodels.MockLoginUser.UserRole,
		})

		user.UpdateUser(c)

		expected := map[string]any{
			"data": map[string]any{
				"name":         *testmodels.MockUpdateUser.Name,
				"username":     *testmodels.MockUpdateUser.Username,
				"email":        testmodels.MockLoginUser.Email,
				"country_code": *testmodels.MockUpdateUser.CountryCode,
				"country":      "Brazil",
				"user_role":    "",
				"role":         *testmodels.MockUpdateUser.Role,
				"interest":     *testmodels.MockUpdateUser.Interest,
				"institution":  *testmodels.MockUpdateUser.Institution,
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
