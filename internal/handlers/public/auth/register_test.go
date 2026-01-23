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
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testutils.SetupTestContext()

	mockCountry := testmodels.NewCountry("", nil)

	mockRegisterUser := testmodels.NewRegisterUser("", "")
	mockUser := models.User{
		Name:        mockRegisterUser.Name,
		Username:    mockRegisterUser.Username,
		Email:       mockRegisterUser.Email,
		Interest:    mockRegisterUser.Interest,
		UserRole:    models.Collaborator,
		Role:        mockRegisterUser.Role,
		CountryID:   mockCountry.ID,
		Country:     mockCountry,
		Institution: mockRegisterUser.Institution,
	}
	lang := "en"
	mockResponse := mockUser.ToResponse(lang)

	t.Run("Success", func(t *testing.T) {
		body := testutils.ToJSON(mockRegisterUser)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)

		svc := &mocks.MockAuthService{
			RegisterFunc: func(ctx context.Context, input models.UserRegisterInput, language string) (*models.UserResponse, error) {
				return &mockResponse, nil
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Register(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "User created successfully. Please wait for activation by an administrator.",
			"data":    mockResponse,
		})

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.RegisterTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/auth/register", tt.Body,
				nil, nil,
			)
			svc := &mocks.MockAuthService{}
			handler := handlerAuth.NewAuthHandler(svc)
			handler.Register(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Email Already Exists", func(t *testing.T) {
		body := testutils.ToJSON(mockRegisterUser)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)

		svc := &mocks.MockAuthService{
			RegisterFunc: func(ctx context.Context, input models.UserRegisterInput, language string) (*models.UserResponse, error) {
				return nil, services.ErrConflictEmail
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Register(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Email is already in use.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Username Already Exists", func(t *testing.T) {
		body := testutils.ToJSON(mockRegisterUser)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)

		svc := &mocks.MockAuthService{
			RegisterFunc: func(ctx context.Context, input models.UserRegisterInput, language string) (*models.UserResponse, error) {
				return nil, services.ErrConflictUsername
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Register(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Username already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Email Mismatch", func(t *testing.T) {
		body := testutils.ToJSON(mockRegisterUser)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)

		svc := &mocks.MockAuthService{
			RegisterFunc: func(ctx context.Context, input models.UserRegisterInput, language string) (*models.UserResponse, error) {
				return nil, services.ErrEmailMismatch
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Register(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Emails must match.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
	t.Run("Error - Passwords Mismatch", func(t *testing.T) {
		body := testutils.ToJSON(mockRegisterUser)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)

		svc := &mocks.MockAuthService{
			RegisterFunc: func(ctx context.Context, input models.UserRegisterInput, language string) (*models.UserResponse, error) {
				return nil, services.ErrPasswordMismatch
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Register(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Passwords must match.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Country Not Found", func(t *testing.T) {
		body := testutils.ToJSON(mockRegisterUser)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)

		svc := &mocks.MockAuthService{
			RegisterFunc: func(ctx context.Context, input models.UserRegisterInput, language string) (*models.UserResponse, error) {
				return nil, services.ErrInvalidCountryCode
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Register(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "No country was found with this code."},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		body := testutils.ToJSON(mockRegisterUser)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/auth/register", body,
			nil, nil,
		)

		svc := &mocks.MockAuthService{
			RegisterFunc: func(ctx context.Context, input models.UserRegisterInput, language string) (*models.UserResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := handlerAuth.NewAuthHandler(svc)
		handler.Register(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again."},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
