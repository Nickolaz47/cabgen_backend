package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestActivateUser(t *testing.T) {
	testutils.SetupTestContext()

	mockAdminUser := testmodels.NewAdminLoginUser()
	mockAdminUserToken := testmodels.NewUserToken(
		mockAdminUser.ID,
		mockAdminUser.Username,
		mockAdminUser.UserRole,
	)

	mockLoginUser := testmodels.NewLoginUser()
	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			ActivateUserFunc: func(ctx context.Context, ID uuid.UUID, adminName string) error {
				return nil
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users/activate", "",
			nil, gin.Params{{Key: "userId", Value: mockLoginUser.ID.String()}},
		)
		c.Set("user", &mockAdminUserToken)
		handler.ActivateUser(c)

		expected := ""

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users/activate", "",
			nil, gin.Params{{Key: "userId", Value: "ifew90843"}},
		)
		c.Set("user", &mockAdminUserToken)
		handler.ActivateUser(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "The URL ID is invalid.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			ActivateUserFunc: func(ctx context.Context, ID uuid.UUID, adminName string) error {
				return services.ErrNotFound
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users/activate", "",
			nil, gin.Params{{Key: "userId", Value: uuid.NewString()}},
		)
		c.Set("user", &mockAdminUserToken)
		handler.ActivateUser(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "User not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			ActivateUserFunc: func(ctx context.Context, ID uuid.UUID, adminName string) error {
				return services.ErrInternal
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users/activate", "",
			nil, gin.Params{{Key: "userId", Value: mockLoginUser.ID.String()}},
		)
		c.Set("user", &mockAdminUserToken)
		handler.ActivateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again."},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
