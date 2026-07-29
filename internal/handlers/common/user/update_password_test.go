package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePassword(t *testing.T) {
	testutils.SetupTestContext()

	mockLoginUser := testmodels.NewLoginUser()
	mockToken := mockLoginUser.ToToken()

	body := `{"current_password":"oldpass","new_password":"newpass123","confirm_password":"newpass123"}`

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockUserService{
			UpdatePasswordFunc: func(ctx context.Context, ID uuid.UUID,
				input models.UpdatePasswordInput) error {
				return nil
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/update-password", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.UpdatePassword(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "Password updated successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Missing user in context", func(t *testing.T) {
		svc := &mocks.MockUserService{}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/update-password", body,
			nil, nil,
		)
		handler.UpdatePassword(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue."},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdatePasswordTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockUserService{}
			handler := user.NewUserHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/users/me/update-password", tt.Body,
				nil, nil,
			)
			c.Set("user", &mockToken)
			handler.UpdatePassword(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Wrong Current Password", func(t *testing.T) {
		svc := &mocks.MockUserService{
			UpdatePasswordFunc: func(ctx context.Context, ID uuid.UUID,
				input models.UpdatePasswordInput) error {
				return services.ErrCurrentPasswordMismatch
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/update-password", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.UpdatePassword(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "The current password is incorrect.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		svc := &mocks.MockUserService{
			UpdatePasswordFunc: func(ctx context.Context, ID uuid.UUID,
				input models.UpdatePasswordInput) error {
				return services.ErrNotFound
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/update-password", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.UpdatePassword(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "User not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
