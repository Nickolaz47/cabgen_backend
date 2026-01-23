package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResolveCookieSecurity(t *testing.T) {
	t.Run("Prod", func(t *testing.T) {
		origEnv := config.Environment
		origHost := config.APIHost
		defer func() {
			config.Environment = origEnv
			config.APIHost = origHost
		}()

		expected := "https://cabgen.com/api"

		config.Environment = "prod"
		config.APIHost = expected

		host, isSecure := resolveCookieSecurity()

		assert.Equal(t, expected, host)
		assert.True(t, isSecure)
	})

	t.Run("Dev", func(t *testing.T) {
		origEnv := config.Environment
		origHost := config.APIHost
		defer func() {
			config.Environment = origEnv
			config.APIHost = origHost
		}()

		expected := "localhost"

		config.Environment = "dev"
		config.APIHost = expected

		host, isSecure := resolveCookieSecurity()

		assert.Equal(t, expected, host)
		assert.False(t, isSecure)
	})
}

func TestCreateCookie(t *testing.T) {
	defer func() {
		os.Unsetenv("API_HOST")
		os.Unsetenv("ENVIRONMENT")
	}()

	host := "localhost"

	os.Setenv("API_HOST", host)
	os.Setenv("ENVIRONMENT", "dev")

	content := "content"
	path := "/"
	expiration := 1 * time.Second

	now := time.Now().Add(time.Second)
	cookie := CreateCookie(
		Access,
		content,
		path,
		expiration,
	)

	assert.Equal(t, Access, cookie.Name)
	assert.Equal(t, content, cookie.Value)
	assert.Equal(t, path, cookie.Path)
	assert.LessOrEqual(t, now, cookie.Expires)
	assert.Equal(t, host, cookie.Domain)
	assert.True(t, cookie.HttpOnly)
	assert.False(t, cookie.Secure)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
}

func TestDeleteCookie(t *testing.T) {
	defer func() {
		os.Unsetenv("API_HOST")
		os.Unsetenv("ENVIRONMENT")
	}()

	host := "localhost"

	os.Setenv("API_HOST", host)
	os.Setenv("ENVIRONMENT", "dev")

	path := "/"

	cookie := DeleteCookie(Access, path)

	assert.Equal(t, Access, cookie.Name)
	assert.Empty(t, cookie.Value)
	assert.Equal(t, path, cookie.Path)
	assert.Equal(t, -1, cookie.MaxAge)
	assert.Equal(t, host, cookie.Domain)
	assert.True(t, cookie.HttpOnly)
	assert.False(t, cookie.Secure)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
}

func TestExtractToken(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(
			http.MethodGet, "/", "", nil, nil,
		)

		c.Request.AddCookie(&http.Cookie{
			Name:     Access,
			Value:    "accessToken",
			Path:     "/",
			Expires:  time.Now().Add(time.Second),
			HttpOnly: true,
		})

		token, err := ExtractToken(c, Access)

		assert.NoError(t, err)
		assert.Equal(t, "accessToken", token)
	})

	t.Run("Error - Missing cookie", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(
			http.MethodGet, "/", "", nil, nil,
		)

		token, err := ExtractToken(c, Access)

		assert.Empty(t, token)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cookie not found")
	})
}

func TestGenerateToken(t *testing.T) {
	provider := NewTokenProvider()

	mockToken := testmodels.NewUserToken(
		uuid.UUID{},
		"nick",
		models.Collaborator,
	)
	secret := []byte("secret")
	expiration := 1 * time.Second

	t.Run("Success", func(t *testing.T) {
		token, err := provider.GenerateToken(mockToken, secret, expiration)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsed := &models.UserToken{}
		tok, err := jwt.ParseWithClaims(token, parsed, func(tk *jwt.Token) (any, error) {
			if tk.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", tk.Header["alg"])
			}
			return secret, nil
		})

		assert.NoError(t, err)
		assert.True(t, tok.Valid)
	})

	t.Run("Error - Invalid expiration time", func(t *testing.T) {
		token, err := provider.GenerateToken(mockToken, secret, -expiration)

		assert.Empty(t, token)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "expiresIn must be > 0")
	})

	t.Run("Error - Empty secret", func(t *testing.T) {
		token, err := provider.GenerateToken(mockToken, nil, expiration)

		assert.Empty(t, token)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "secret is empty")
	})
}

func TestValidateToken(t *testing.T) {
	provider := NewTokenProvider()
	secret := []byte("access_secret")

	mockUserToken := testmodels.NewUserToken(
		uuid.UUID{}, "nick", models.Collaborator,
	)

	t.Run("Success", func(t *testing.T) {
		tokenStr, err := provider.GenerateToken(
			mockUserToken,
			secret,
			time.Second,
		)
		assert.NoError(t, err)

		result, err := provider.ValidateToken(tokenStr, secret)

		assert.NoError(t, err)
		assert.Equal(t, mockUserToken.ID, result.ID)
		assert.Equal(t, mockUserToken.Username, result.Username)
		assert.Equal(t, mockUserToken.UserRole, result.UserRole)
	})

	t.Run("Error - Invalid signing method", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, mockUserToken)
		tokenStr, _ := token.SignedString(secret)

		parts := strings.Split(tokenStr, ".")
		header := `{"alg":"RS256","typ":"JWT"}`
		encodedHeader := base64.RawURLEncoding.EncodeToString([]byte(header))
		tokenStr = encodedHeader + "." + parts[1] + "." + parts[2]

		result, err := provider.ValidateToken(tokenStr, secret)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected signing method")
	})

	t.Run("Error - Token expired", func(t *testing.T) {
		expiredClaims := mockUserToken
		expiredClaims.RegisteredClaims.ExpiresAt =
			jwt.NewNumericDate(time.Now().Add(-time.Hour))

		expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).
			SignedString(secret)

		result, err := provider.ValidateToken(expiredToken, secret)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token expired")
	})

	t.Run("Error - Invalid token", func(t *testing.T) {
		result, err := provider.ValidateToken("not.a.valid.token", secret)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("Error - Invalid or expired token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"ID":       []string{"invalid-uuid"},
			"Username": 12345,
			"UserRole": true,
		})
		tokenStr, _ := token.SignedString(secret)

		result, err := provider.ValidateToken(tokenStr, secret)

		assert.Nil(t, result)
		assert.EqualError(t, err, "invalid or expired token")
	})
}

func TestGetSecretKey(t *testing.T) {
	origAccessKey := config.AccessKey
	origRefreshKey := config.RefreshKey
	defer func() {
		config.AccessKey = origAccessKey
		config.RefreshKey = origRefreshKey
	}()

	accessSecret := []byte("access_secret")
	refreshSecret := []byte("refresh_secret")
	config.AccessKey = accessSecret
	config.RefreshKey = refreshSecret

	t.Run("Access Cookie", func(t *testing.T) {
		secret, err := GetSecretKey(Access)

		assert.NoError(t, err)
		assert.Equal(t, []byte(accessSecret), secret)
	})

	t.Run("Refresh Cookie", func(t *testing.T) {
		secret, err := GetSecretKey(Refresh)

		assert.NoError(t, err)
		assert.Equal(t, []byte(refreshSecret), secret)
	})

	t.Run("Invalid Cookie", func(t *testing.T) {
		var Invalid Cookie = "Invalid"

		secret, err := GetSecretKey(Invalid)

		assert.Empty(t, secret)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid cookie name")
	})
}
