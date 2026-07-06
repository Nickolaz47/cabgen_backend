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

func TestForgotPassword(t *testing.T) {
	testutils.SetupTestContext()

	mockInput := models.ForgotPasswordInput{
		Email: "test@mail.com",
	}

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/forgot-password",
			testutils.ToJSON(mockInput),
			nil, nil,
		)
		svc := &mocks.MockAuthService{
			ForgotPasswordFunc: func(ctx context.Context,
				input models.ForgotPasswordInput) error {
				return nil
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.ForgotPassword(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "If the email is registered, you will receive" +
					" password reset instructions shortly."},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.ForgotPasswordBadRequestTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/auth/forgot-password", tt.Body,
				nil, nil,
			)
			svc := &mocks.MockAuthService{}
			handler := handlerAuth.NewAuthHandler(svc)
			handler.ForgotPassword(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Internal Server", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/forgot-password",
			testutils.ToJSON(mockInput), nil, nil,
		)
		svc := &mocks.MockAuthService{
			ForgotPasswordFunc: func(ctx context.Context,
				input models.ForgotPasswordInput) error {
				return services.ErrInternal
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.ForgotPassword(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again."},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
