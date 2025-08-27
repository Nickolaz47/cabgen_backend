package public_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	testutils.SetupTestContext()

	db := testutils.SetupTestRepos()
	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	t.Run(data.LoginSuccess.Name, func(t *testing.T) {
		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", data.LoginSuccess.Body)
		public.Login(c)

		cookies := w.Result().Cookies()
		var accessCookie, refreshCookie string

		for _, cookie := range cookies {
			if cookie.Name == auth.Access {
				accessCookie = cookie.Value
			}
			if cookie.Name == auth.Refresh {
				refreshCookie = cookie.Value
			}
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, data.LoginSuccess.Expected, w.Body.String())
		assert.NotEmpty(t, accessCookie)
		assert.NotEmpty(t, refreshCookie)
	})

	for _, tt := range data.LoginBadRequestTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", tt.Body)
			public.Login(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	for _, tt := range data.LoginUnauthorizedTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", tt.Body)
			public.Login(c)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	db.Model(&models.User{}).
		Where("username = ?", mockLoginUser.Username).
		Update("is_active", false)

	t.Run(data.LoginUserDeactivatedTest.Name, func(t *testing.T) {
		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", data.LoginUserDeactivatedTest.Body)
		public.Login(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t, data.LoginUserDeactivatedTest.Expected, w.Body.String())
	})
}
