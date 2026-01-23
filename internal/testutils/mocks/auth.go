package mocks

import (
	"context"
	"net/http"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/gin-gonic/gin"
)

type MockHasher struct {
	HashFunc          func(password string) (string, error)
	CheckPasswordFunc func(hashPassword, password string) error
}

func (h *MockHasher) Hash(password string) (string, error) {
	if h.HashFunc != nil {
		return h.HashFunc(password)
	}
	return "", nil
}

func (h *MockHasher) CheckPassword(hashPassword, password string) error {
	if h.CheckPasswordFunc != nil {
		return h.CheckPasswordFunc(hashPassword, password)
	}
	return nil
}

type MockTokenProvider struct {
	ExtractTokenFunc  func(c *gin.Context, cookieName auth.Cookie) (string, error)
	GenerateTokenFunc func(user models.UserToken, secret []byte, expiresIn time.Duration) (string, error)
	ValidateTokenFunc func(tokenStr string, secret []byte) (*models.UserToken, error)
}

func (p *MockTokenProvider) ExtractToken(
	c *gin.Context, cookieName auth.Cookie) (string, error) {
	if p.ExtractTokenFunc != nil {
		return p.ExtractTokenFunc(c, cookieName)
	}
	return "", nil
}

func (p *MockTokenProvider) GenerateToken(
	user models.UserToken, secret []byte,
	expiresIn time.Duration) (string, error) {
	if p.GenerateTokenFunc != nil {
		return p.GenerateTokenFunc(user, secret, expiresIn)
	}
	return "", nil
}

func (p *MockTokenProvider) ValidateToken(
	tokenStr string, secret []byte) (*models.UserToken, error) {
	if p.ValidateTokenFunc != nil {
		return p.ValidateTokenFunc(tokenStr, secret)
	}
	return nil, nil
}

type MockAuthService struct {
	RegisterFunc func(
		ctx context.Context,
		input models.UserRegisterInput,
		language string) (*models.UserResponse, error)
	LoginFunc func(ctx context.Context,
		input models.LoginInput) (*models.Cookies, error)
	RefreshFunc func(ctx context.Context,
		tokenStr string) (*http.Cookie, error)
}

func (s *MockAuthService) Register(
	ctx context.Context,
	input models.UserRegisterInput,
	language string) (*models.UserResponse, error) {
	if s.RegisterFunc != nil {
		return s.RegisterFunc(ctx, input, language)
	}
	return nil, nil
}

func (s *MockAuthService) Login(ctx context.Context,
	input models.LoginInput) (*models.Cookies, error) {
	if s.LoginFunc != nil {
		return s.LoginFunc(ctx, input)
	}
	return nil, nil
}

func (s *MockAuthService) Refresh(ctx context.Context,
	tokenStr string) (*http.Cookie, error) {
	if s.RefreshFunc != nil {
		return s.RefreshFunc(ctx, tokenStr)
	}
	return nil, nil
}
