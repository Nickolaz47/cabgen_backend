package auth_test

import (
	"context"
	"net/http"
	"testing"

	handlerAuth "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestResetPassword(t *testing.T) {
	testutils.SetupTestContext()

	mockInput := models.ResetPasswordInput{
		Token:           "valid-uuid-token-string",
		NewPassword:     "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/reset-password",
			testutils.ToJSON(mockInput),
			nil, nil,
		)
		svc := &mocks.MockAuthService{
			ResetPasswordFunc: func(ctx context.Context,
				input models.ResetPasswordInput) error {
				return nil
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.ResetPassword(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "Password reset successfully. Please log in."},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.ResetPasswordBadRequestTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/auth/reset-password", tt.Body,
				nil, nil,
			)
			svc := &mocks.MockAuthService{}
			handler := handlerAuth.NewAuthHandler(svc)
			handler.ResetPassword(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Invalid Token", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/reset-password",
			testutils.ToJSON(mockInput), nil, nil,
		)
		svc := &mocks.MockAuthService{
			ResetPasswordFunc: func(ctx context.Context,
				input models.ResetPasswordInput) error {
				return services.ErrInvalidToken
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.ResetPassword(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Invalid password reset link. Please request" +
					" a new one."},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Expired Token", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/reset-password",
			testutils.ToJSON(mockInput), nil, nil,
		)
		svc := &mocks.MockAuthService{
			ResetPasswordFunc: func(ctx context.Context,
				input models.ResetPasswordInput) error {
				return services.ErrExpiredToken
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.ResetPassword(c)

		expected := testutils.ToJSON(
			map[string]string{"error": "Password reset link expired. Please" +
				" request a new one."},
		)

		assert.Equal(t, http.StatusGone, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/reset-password",
			testutils.ToJSON(mockInput), nil, nil,
		)
		svc := &mocks.MockAuthService{
			ResetPasswordFunc: func(ctx context.Context,
				input models.ResetPasswordInput) error {
				return services.ErrInternal
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.ResetPassword(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again."},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
