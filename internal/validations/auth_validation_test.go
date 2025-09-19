package validations_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserTokenFromContext(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		mockToken := testmodels.NewUserToken(
			uuid.UUID{}, "nick", models.Collaborator,
		)

		c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)
		c.Set("user", &mockToken)

		token, ok := validations.GetUserTokenFromContext(c)

		assert.True(t, ok)
		assert.Equal(t, &mockToken, token)
	})

	t.Run("Invalid token type", func(t *testing.T) {
		invalidToken := models.Admin

		c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)
		c.Set("user", &invalidToken)

		token, ok := validations.GetUserTokenFromContext(c)

		assert.False(t, ok)
		assert.Empty(t, token)
	})

	t.Run("Token missing", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)
		token, ok := validations.GetUserTokenFromContext(c)

		assert.False(t, ok)
		assert.Empty(t, token)
	})
}
