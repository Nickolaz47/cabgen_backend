package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/gin-gonic/gin"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Cookie string

const (
	Access                 = "AccessCookie"
	Refresh                = "RefreshCookie"
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour
)

func (c Cookie) String() string {
	return string(c)
}

func CreateCookie(cookieName Cookie, cookieContent, path string, expiration time.Duration) *http.Cookie {
	host := "localhost"
	var secure bool
	if config.Environment == "prod" {
		secure = true
		host = config.APIHost
	}

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

func GenerateToken(user models.UserToken, secret []byte, expiresIn time.Duration) (string, error) {
	user.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(secret)
}

func ValidateToken(c *gin.Context, cookieName Cookie, secret []byte) (*models.UserToken, error) {
	tokenStr, err := extractToken(c, cookieName)
	if err != nil {
		return nil, fmt.Errorf("cookie not found: %v", err)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &models.UserToken{}, func(token *jwt.Token) (any, error) {
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

	if claims, ok := token.Claims.(*models.UserToken); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid or expired token")
}

func extractToken(c *gin.Context, cookieName Cookie) (string, error) {
	token, err := c.Cookie(cookieName.String())
	if err != nil {
		return "", fmt.Errorf("cookie not found: %v", err)
	}

	return token, nil
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
