package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/user"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	testutils.SetupTestContext()

	mockLoginUser := testmodels.NewLoginUser()
	mockToken := mockLoginUser.ToToken()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockUserService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return nil
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/users/me", "",
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.DeleteUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "User deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Missing user in context", func(t *testing.T) {
		svc := &mocks.MockUserService{}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/users/me", "",
			nil, nil,
		)
		handler.DeleteUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue."},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockUserService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrNotFound
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/users/me", "",
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.DeleteUser(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "User not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockUserService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrInternal
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/users/me", "",
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.DeleteUser(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
