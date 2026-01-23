package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetOwnUser(t *testing.T) {
	testutils.SetupTestContext()
	lang := "en"

	mockLoginUser := testmodels.NewLoginUser()
	mockToken := mockLoginUser.ToToken()

	response := mockLoginUser.ToResponse(lang)

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error) {
				return &response, nil
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/users/me", "",
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.GetOwnUser(c)

		expected := testutils.ToJSON(
			map[string]models.UserResponse{
				"data": response,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Missing user in context", func(t *testing.T) {
		svc := &mocks.MockUserService{}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/users/me", "",
			nil, nil,
		)
		handler.GetOwnUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue."},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/users/me", "",
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.GetOwnUser(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "User not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/users/me", "",
			nil, nil,
		)
		c.Set("user", &mockToken)
		handler.GetOwnUser(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
