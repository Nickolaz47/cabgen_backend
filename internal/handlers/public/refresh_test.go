package public_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	db := testutils.SetupTestRepos()
	origAccessKey := config.AccessKey
	origRefreshKey := config.RefreshKey
	config.AccessKey = []byte("access_secret")
	config.RefreshKey = []byte("refresh_secret")

	defer func() {
		config.AccessKey = origAccessKey
		config.RefreshKey = origRefreshKey
	}()

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	mockLoginInput := models.LoginInput{
		Username: mockLoginUser.Username,
		Password: "12345678",
	}

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/auth/login",
			testutils.ToJSON(mockLoginInput),
			nil, nil,
		)

		public.Login(c)
		var accessCookie, refreshCookie string

		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == auth.Refresh {
				refreshCookie = cookie.Value
			}
		}

		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(nil))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: ""})
		req.AddCookie(&http.Cookie{Name: auth.Refresh, Value: refreshCookie})

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		public.Refresh(c)

		cookies = w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == auth.Access {
				accessCookie = cookie.Value
			}
		}

		expected := `{"message": "Access token renewed successfully."}`

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		assert.NotEmpty(t, accessCookie)
	})

	t.Run("Invalid token", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/auth/login",
			testutils.ToJSON(mockLoginInput),
			nil, nil,
		)

		public.Login(c)

		var refreshCookie string

		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == auth.Refresh {
				refreshCookie = cookie.Value
			}
		}

		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(nil))
		req.AddCookie(&http.Cookie{Name: auth.Refresh, Value: refreshCookie[:len(refreshCookie)-5]})

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		public.Refresh(c)

		expected := `{"error": "Unauthorized. Please log in to continue."}`

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
