package user_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	testutils.SetupTestContext()

	db := testutils.SetupTestRepos()

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/user/nick", "",
			nil, gin.Params{{Key: "username", Value: "nick"}},
		)

		user.DeleteUser(c)

		expected := `{"message": "User deleted successfully."}`

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("User not found", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/user/xxxx", "",
			nil, gin.Params{{Key: "username", Value: "xxxx"}},
		)

		user.DeleteUser(c)

		expected := `{"error": "User not found."}`

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
