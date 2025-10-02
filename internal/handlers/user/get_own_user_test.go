package user_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetOwnUser(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/auth/login",
			testutils.ToJSON(models.LoginInput{
				Username: mockLoginUser.Username,
				Password: "12345678"},
			), nil, nil,
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
			ID:       mockLoginUser.ID,
			Username: mockLoginUser.Username,
			UserRole: mockLoginUser.UserRole,
		})

		user.GetOwnUser(c)

		expected := testutils.ToJSON(responses.APIResponse{Data: mockLoginUser.ToResponse(c)})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Missing user in context", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/user/me",
			"", nil, nil,
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
			"", nil, nil,
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
