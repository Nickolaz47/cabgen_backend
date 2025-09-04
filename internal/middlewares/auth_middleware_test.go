package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func writeMockEnvFile(t *testing.T, envFilePath, envContent string) {
	if err := os.WriteFile(envFilePath, []byte(envContent), 0644); err != nil {
		t.Errorf("failed to write mock env file: %v", err)
	}
}

func TestAuthMiddleware(t *testing.T) {
	testutils.SetupTestContext()

	mockLoginUser := testmodels.NewLoginUser()
	mockToken := testmodels.NewUserToken(
		mockLoginUser.ID,
		mockLoginUser.Username,
		mockLoginUser.UserRole,
	)

	tempDir := t.TempDir()
	mockEnvFile := filepath.Join(tempDir, "test.env")
	mockEnvContent := `
	SECRET_ACCESS_KEY=super-test-secret
	PORT=8080
	`

	writeMockEnvFile(t, mockEnvFile, mockEnvContent)
	config.LoadEnvVariables(mockEnvFile)

	secret, _ := auth.GetSecretKey(auth.Access)

	t.Run("Success", func(t *testing.T) {
		w, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.AuthMiddleware())

		r.GET("/", func(c *gin.Context) {
			rawUserToken, exists := c.Get("user")
			if !exists {
				c.JSON(http.StatusInternalServerError, map[string]*models.UserToken{
					"userToken": nil,
				})
				return
			}

			userToken, ok := rawUserToken.(*models.UserToken)
			if !ok {
				c.JSON(http.StatusInternalServerError, map[string]*models.UserToken{
					"userToken": nil,
				})
				return
			}

			c.JSON(http.StatusOK, map[string]*models.UserToken{
				"userToken": userToken,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		mockAccessToken, _ := auth.GenerateToken(
			mockToken, secret, auth.AccessTokenExpiration,
		)
		mockAccessCookie := auth.CreateCookie(
			auth.Access, mockAccessToken,
			"/", auth.AccessTokenExpiration,
		)

		req.AddCookie(mockAccessCookie)
		r.ServeHTTP(w, req)

		expected := testutils.ToJSON(map[string]models.UserToken{"userToken": mockToken})

		var got map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)

		if data, ok := got["userToken"].(map[string]any); ok {
			delete(data, "sub")
			delete(data, "exp")
			delete(data, "iat")
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, testutils.ToJSON(got))
	})

	t.Run("Expired token", func(t *testing.T) {
		w, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.AuthMiddleware())

		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		mockAccessToken, _ := auth.GenerateToken(
			mockToken, secret, -1*time.Second,
		)
		mockAccessCookie := auth.CreateCookie(
			auth.Access, mockAccessToken,
			"/", -1*time.Second,
		)

		req.AddCookie(mockAccessCookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Any token error", func(t *testing.T) {
		w, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.AuthMiddleware())

		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		mockAccessToken, _ := auth.GenerateToken(
			mockToken, secret[:10], auth.AccessTokenExpiration,
		)
		mockAccessCookie := auth.CreateCookie(
			auth.Access, mockAccessToken,
			"/", auth.AccessTokenExpiration,
		)

		req.AddCookie(mockAccessCookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
