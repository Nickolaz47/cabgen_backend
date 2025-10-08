package user_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserActivation(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockAdminLoginUser := testmodels.NewAdminLoginUser()
	db.Create(&mockAdminLoginUser)

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	mockAdminUserToken := testmodels.NewUserToken(
		mockAdminLoginUser.ID,
		mockAdminLoginUser.Username,
		mockAdminLoginUser.UserRole,
	)

	gc, gw := testutils.SetupGinContext(
		http.MethodPost,
		"/api/auth/login",
		testutils.ToJSON(models.LoginInput{
			Username: mockAdminLoginUser.Username,
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

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/admin/user/nick", bytes.NewBuffer(nil))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "nick"}}

		c.Set("user", &mockAdminUserToken)

		user.UpdateUserActivation(c)

		var updatedUser models.User
		db.Where(`username = ?`, mockLoginUser.Username).First(&updatedUser)

		expected := `{"message": "User deactivated successfully."}`

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		assert.Equal(t, mockAdminLoginUser.Username, *updatedUser.ActivatedBy)
		assert.NotEmpty(t, updatedUser.ActivatedOn)
	})

	t.Run("Missing user in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/admin/user/nick", bytes.NewBuffer(nil))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "nick"}}

		user.UpdateUserActivation(c)

		expected := `{"error": "Unauthorized. Please log in to continue."}`

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("User not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/admin/user/nick", bytes.NewBuffer(nil))
		req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "username", Value: "xxxx"}}

		c.Set("user", &mockAdminUserToken)

		user.UpdateUserActivation(c)

		expected := `{"error": "User not found."}`

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
