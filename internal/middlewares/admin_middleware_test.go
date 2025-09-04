package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"

	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminMiddleware(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		w, r := testutils.SetupMiddlewareContext()

		mockAuthMiddleware := func(c *gin.Context) {
			mockUserToken := testmodels.NewUserToken(
				uuid.UUID{}, "nick", models.Admin,
			)
			c.Set("user", &mockUserToken)
		}

		testutils.AddMiddlewares(r, mockAuthMiddleware, middlewares.AdminMiddleware())
		testutils.AddTestGetRoute(r, http.StatusOK)
		testutils.DoGetRequest(r, w)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing user token", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := gin.New()

		testutils.AddMiddlewares(r, middlewares.AdminMiddleware())
		testutils.AddTestGetRoute(r, http.StatusOK)
		testutils.DoGetRequest(r, w)

		expected := `{"error":"Unauthorized. Please log in to continue."}`

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("User is not admin", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := gin.New()

		mockAuthMiddleware := func(c *gin.Context) {
			mockUserToken := testmodels.NewUserToken(
				uuid.UUID{}, "nick", models.Collaborator,
			)
			c.Set("user", &mockUserToken)
		}

		testutils.AddMiddlewares(r, mockAuthMiddleware, middlewares.AdminMiddleware())
		testutils.AddTestGetRoute(r, http.StatusOK)
		testutils.DoGetRequest(r, w)

		expected := `{"error":"Unauthorized. Please log in to continue."}`

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
