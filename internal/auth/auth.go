package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/gin-gonic/gin"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	Access                 = "AccessCookie"
	Refresh                = "RefreshCookie"
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour
)

type Cookie string

func (c Cookie) String() string {
	return string(c)
}

func resolveCookieSecurity() (host string, secure bool) {
	if config.Environment == "prod" {
		return config.APIHost, true
	}
	return "localhost", false
}

func GetSecretKey(cookieName Cookie) ([]byte, error) {
	var secret []byte
	switch cookieName {
	case Access:
		secret = config.AccessKey
	case Refresh:
		secret = config.RefreshKey
	default:
		return nil, errors.New("invalid cookie name")
	}

	return secret, nil
}

func CreateCookie(cookieName Cookie, cookieContent, path string, expiration time.Duration) *http.Cookie {
	host, secure := resolveCookieSecurity()

	return &http.Cookie{
		Name:     cookieName.String(),
		Value:    cookieContent,
		Path:     path,
		Domain:   host,
		Expires:  time.Now().Add(expiration),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func DeleteCookie(cookieName Cookie, path string) *http.Cookie {
	host, secure := resolveCookieSecurity()

	return &http.Cookie{
		Name:     cookieName.String(),
		Value:    "",
		Path:     path,
		Domain:   host,
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func ExtractToken(c *gin.Context, cookieName Cookie) (string, error) {
	token, err := c.Cookie(cookieName.String())
	if err != nil {
		return "", fmt.Errorf("cookie not found: %v", err)
	}

	return token, nil
}

type TokenProvider interface {
	GenerateToken(user models.UserToken, secret []byte, expiresIn time.Duration) (string, error)
	ValidateToken(tokenStr string, secret []byte) (*models.UserToken, error)
}

type tokenProvider struct{}

func NewTokenProvider() TokenProvider {
	return &tokenProvider{}
}

func (t *tokenProvider) GenerateToken(user models.UserToken, secret []byte, expiresIn time.Duration) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf("secret is empty")
	}
	if expiresIn <= 0 {
		return "", fmt.Errorf("expiresIn must be > 0")
	}

	user.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(secret)
}

func (t *tokenProvider) ValidateToken(tokenStr string, secret []byte) (*models.UserToken, error) {
	// Parses the token using MapClaims
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired: %v", err)
		}
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Validates and converts claims to UserToken
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data, err := json.Marshal(claims)
		if err != nil {
			return nil, errors.New("invalid or expired token")
		}

		var userToken models.UserToken
		if err := json.Unmarshal(data, &userToken); err != nil {
			return nil, errors.New("invalid or expired token")
		}

		return &userToken, nil
	}

	return nil, errors.New("invalid or expired token")
}
