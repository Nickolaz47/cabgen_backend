package public_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	db := testutils.SetupTestRepos()

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	c, w := testutils.SetupGinContext(
		http.MethodPost,
		"/api/auth/login",
		testutils.ToJSON(models.LoginInput{
			Username: mockLoginUser.Username,
			Password: "12345678"},
		), nil, nil,
	)

	public.Login(c)
	var accessCookie, refreshCookie string

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == auth.Access {
			accessCookie = cookie.Value
		}
		if cookie.Name == auth.Refresh {
			refreshCookie = cookie.Value
		}
	}

	assert.NotEmpty(t, accessCookie)
	assert.NotEmpty(t, refreshCookie)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBuffer(nil))
	req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})
	req.AddCookie(&http.Cookie{Name: auth.Refresh, Value: refreshCookie})

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req

	public.Logout(c)

	cookies = w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == auth.Access {
			accessCookie = cookie.Value
		}
		if cookie.Name == auth.Refresh {
			refreshCookie = cookie.Value
		}
	}

	expected := `{"message": "Logout successful."}`

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
	assert.Empty(t, accessCookie)
	assert.Empty(t, refreshCookie)
}
