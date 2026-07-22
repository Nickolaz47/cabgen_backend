package auth_test

import (
	"net/http"
	"testing"

	authHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/common/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestMe(t *testing.T) {
	testutils.SetupTestContext()

	mockLoginUser := testmodels.NewLoginUser()
	mockToken := mockLoginUser.ToToken()

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/auth/me", "", nil, nil,
		)
		c.Set("user", &mockToken)

		handler := authHandler.NewAuthHandler(nil)
		handler.Me(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": map[string]any{
					"id":        mockToken.ID,
					"username":  mockToken.Username,
					"user_role": mockToken.UserRole,
				},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/auth/me", "", nil, nil,
		)

		handler := authHandler.NewAuthHandler(nil)
		handler.Me(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue."},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
