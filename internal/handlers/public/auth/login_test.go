package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	handlerAuth "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	testutils.SetupTestContext()

	accessCookie := &http.Cookie{
		Name:     auth.Access,
		Value:    "accessToken",
		Path:     "/",
		Expires:  time.Now().Add(time.Second),
		HttpOnly: true,
	}
	refreshCookie := &http.Cookie{
		Name:     auth.Refresh,
		Value:    "refreshToken",
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}

	mockLogin := models.LoginInput{
		Username: "nikol47",
		Password: "12345678",
	}

	t.Run(data.LoginSuccess.Name, func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/login",
			data.LoginSuccess.Body,
			nil, nil,
		)
		svc := &mocks.MockAuthService{
			LoginFunc: func(ctx context.Context, input models.LoginInput) (*models.Cookies, error) {
				return &models.Cookies{
					AccessCookie:  accessCookie,
					RefreshCookie: refreshCookie,
				}, nil
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Login(c)

		cookies := w.Result().Cookies()
		var accessCookie, refreshCookie *http.Cookie

		for _, cookie := range cookies {
			if cookie.Name == auth.Access {
				accessCookie = cookie
			}
			if cookie.Name == auth.Refresh {
				refreshCookie = cookie
			}
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, data.LoginSuccess.Expected, w.Body.String())

		assert.NotEmpty(t, accessCookie)
		assert.NotEmpty(t, refreshCookie)
	})

	for _, tt := range data.LoginBadRequestTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/auth/login", tt.Body,
				nil, nil,
			)
			svc := &mocks.MockAuthService{}
			handler := handlerAuth.NewAuthHandler(svc)
			handler.Login(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Invalid Credentials", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/login",
			testutils.ToJSON(mockLogin), nil, nil,
		)
		svc := &mocks.MockAuthService{
			LoginFunc: func(ctx context.Context, input models.LoginInput) (*models.Cookies, error) {
				return nil, services.ErrInvalidCredentials
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Login(c)

		expected := testutils.ToJSON(
			map[string]string{"error": "Invalid credentials."},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Inactive User", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/login",
			testutils.ToJSON(mockLogin), nil, nil,
		)
		svc := &mocks.MockAuthService{
			LoginFunc: func(ctx context.Context, input models.LoginInput) (*models.Cookies, error) {
				return nil, services.ErrDisabledUser
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Login(c)

		expected := testutils.ToJSON(
			map[string]string{"error": "Your account is not activated."},
		)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/login",
			testutils.ToJSON(mockLogin), nil, nil,
		)
		svc := &mocks.MockAuthService{
			LoginFunc: func(ctx context.Context, input models.LoginInput) (*models.Cookies, error) {
				return nil, services.ErrInternal
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Login(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again."},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
