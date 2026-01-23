package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	authHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	testutils.SetupTestContext()

	accessCookie := &http.Cookie{
		Name:     auth.Access,
		Value:    "accessToken",
		Path:     "/",
		Expires:  time.Now().Add(time.Second),
		HttpOnly: true,
	}
	refreshCookie := &http.Cookie{
		Name:     auth.Refresh,
		Value:    "refreshToken",
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}

	c, w := testutils.SetupGinContext(
		http.MethodPost, "/api/auth/logout", "", nil, nil,
	)
	c.Request.AddCookie(accessCookie)
	c.Request.AddCookie(refreshCookie)

	handler := authHandler.NewAuthHandler(nil)
	handler.Logout(c)

	var responseAccess, responseRefresh *http.Cookie
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == auth.Access {
			responseAccess = cookie
		}
		if cookie.Name == auth.Refresh {
			responseRefresh = cookie
		}
	}

	expected := testutils.ToJSON(
		map[string]string{"message": "Logout successful."},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())

	assert.Empty(t, responseAccess.Value)
	assert.Less(t, responseAccess.Expires, time.Now())
	assert.Negative(t, responseAccess.MaxAge)

	assert.Empty(t, responseRefresh.Value)
	assert.Less(t, responseRefresh.Expires, time.Now())
	assert.Negative(t, responseRefresh.MaxAge)
}
