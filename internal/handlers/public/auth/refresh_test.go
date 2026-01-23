package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	authHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	testutils.SetupTestContext()

	accessCookie := &http.Cookie{
		Name:     auth.Access,
		Value:    "accessToken",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	refreshCookie := &http.Cookie{
		Name:     auth.Refresh,
		Value:    "refreshToken",
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}

	// mockLoginUser := testmodels.NewLoginUser()
	// lang := "en"
	// mockResponse := mockLoginUser.ToResponse(lang)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/auth/refresh",
			"", nil, nil,
		)
		c.Request.AddCookie(accessCookie)
		c.Request.AddCookie(refreshCookie)

		svc := &mocks.MockAuthService{
			RefreshFunc: func(
				ctx context.Context, tokenStr string) (*http.Cookie, error) {
				newCookie := accessCookie
				newCookie.Expires = time.Now().Add(15 * time.Minute)
				return newCookie, nil
			},
		}
		handler := authHandler.NewAuthHandler(svc)
		handler.Refresh(c)

		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == auth.Access {
				accessCookie = cookie
			}
		}

		expected := testutils.ToJSON(
			map[string]string{"message": "Access token renewed successfully."})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())

		assert.NotEmpty(t, accessCookie.Value)
		assert.Greater(t, accessCookie.Expires, time.Now())
	})

	t.Run("Error - Missing Refresh Token", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/auth/refresh",
			"", nil, nil,
		)

		svc := &mocks.MockAuthService{}
		handler := authHandler.NewAuthHandler(svc)
		handler.Refresh(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue."},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/auth/refresh",
			"", nil, nil,
		)

		svc := &mocks.MockAuthService{
			RefreshFunc: func(ctx context.Context, tokenStr string) (*http.Cookie, error) {
				return nil, services.ErrUnauthorized
			},
		}
		handler := authHandler.NewAuthHandler(svc)
		handler.Refresh(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue."},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/auth/refresh",
			"", nil, nil,
		)
		c.Request.AddCookie(accessCookie)
		c.Request.AddCookie(refreshCookie)

		svc := &mocks.MockAuthService{
			RefreshFunc: func(ctx context.Context, tokenStr string) (*http.Cookie, error) {
				return nil, services.ErrInternal
			},
		}
		handler := authHandler.NewAuthHandler(svc)
		handler.Refresh(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
