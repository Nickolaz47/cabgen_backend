package user_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByUsername(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/user/nick",
			"", nil, gin.Params{{Key: "username", Value: "nick"}},
		)

		user.GetUserByUsername(c)

		expected := testutils.ToJSON(
			map[string]models.User{
				"data": mockLoginUser,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("User not found", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/user/jao",
			"", nil, gin.Params{{Key: "username", Value: "jao"}},
		)

		user.GetUserByUsername(c)

		expected := `{"error": "User not found."}`

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
