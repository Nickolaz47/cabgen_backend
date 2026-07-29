package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestEmailUpdate(t *testing.T) {
	testutils.SetupTestContext()

	mockLoginUser := testmodels.NewLoginUser()
	mockToken := mockLoginUser.ToToken()

	body := `{"new_email":"new@example.com","confirm_new_email":"new@example.com"}`

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockUserService{
			RequestEmailUpdateFunc: func(ctx context.Context, ID uuid.UUID,
				input models.RequestEmailUpdateInput) error {
				return nil
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/request-email-update", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.RequestEmailUpdate(c)

		expected := testutils.ToJSON(map[string]string{
			"message": "A confirmation link has been sent to your new email address. Please check your inbox.",
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Missing user in context", func(t *testing.T) {
		svc := &mocks.MockUserService{}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/request-email-update", body,
			nil, nil,
		)
		handler.RequestEmailUpdate(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Unauthorized. Please log in to continue.",
		})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Emails do not match", func(t *testing.T) {
		svc := &mocks.MockUserService{}
		handler := user.NewUserHandler(svc)

		invalidBody := `{"new_email":"new@example.com","confirm_new_email":"other@example.com"}`
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/request-email-update", invalidBody,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.RequestEmailUpdate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Email same as current", func(t *testing.T) {
		svc := &mocks.MockUserService{
			RequestEmailUpdateFunc: func(ctx context.Context, ID uuid.UUID,
				input models.RequestEmailUpdateInput) error {
				return services.ErrEmailSame
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/request-email-update", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.RequestEmailUpdate(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "The new email is the same as your current email.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		svc := &mocks.MockUserService{
			RequestEmailUpdateFunc: func(ctx context.Context, ID uuid.UUID,
				input models.RequestEmailUpdateInput) error {
				return services.ErrNotFound
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/request-email-update", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.RequestEmailUpdate(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "User not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}

func TestConfirmEmailUpdate(t *testing.T) {
	testutils.SetupTestContext()

	mockLoginUser := testmodels.NewLoginUser()
	mockToken := mockLoginUser.ToToken()

	body := `{"token":"valid-token"}`

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockUserService{
			ConfirmEmailUpdateFunc: func(ctx context.Context, ID uuid.UUID,
				input models.ConfirmEmailUpdateInput) error {
				return nil
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/confirm-email-update", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.ConfirmEmailUpdate(c)

		expected := testutils.ToJSON(map[string]string{
			"message": "Email updated successfully.",
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Missing user in context", func(t *testing.T) {
		svc := &mocks.MockUserService{}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/confirm-email-update", body,
			nil, nil,
		)
		handler.ConfirmEmailUpdate(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Unauthorized. Please log in to continue.",
		})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Token", func(t *testing.T) {
		svc := &mocks.MockUserService{
			ConfirmEmailUpdateFunc: func(ctx context.Context, ID uuid.UUID,
				input models.ConfirmEmailUpdateInput) error {
				return services.ErrInvalidEmailUpdateToken
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/confirm-email-update", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.ConfirmEmailUpdate(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Invalid email update link. Please request a new one.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Expired Token", func(t *testing.T) {
		svc := &mocks.MockUserService{
			ConfirmEmailUpdateFunc: func(ctx context.Context, ID uuid.UUID,
				input models.ConfirmEmailUpdateInput) error {
				return services.ErrExpiredEmailUpdateToken
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/users/me/confirm-email-update", body,
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.ConfirmEmailUpdate(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Email update link expired. Please request a new one.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
