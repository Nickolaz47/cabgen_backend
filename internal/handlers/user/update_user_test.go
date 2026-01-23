package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	testutils.SetupTestContext()

	mockUser := testmodels.NewLoginUser()
	mockUserToken := testmodels.NewUserToken(
		mockUser.ID,
		mockUser.Username,
		mockUser.UserRole,
	)

	updateInput := testmodels.NewUserUpdateInput()
	updateResponse := models.UserResponse{
		Name:        *updateInput.Name,
		Username:    *updateInput.Username,
		Email:       mockUser.Email,
		CountryCode: *updateInput.CountryCode,
		Country:     "Brazil",
		UserRole:    mockUser.UserRole,
		Role:        updateInput.Role,
		Interest:    updateInput.Interest,
		Institution: updateInput.Institution,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockUserService{
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.UserUpdateInput,
				language string,
			) (*models.UserResponse, error) {
				return &updateResponse, nil
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/users",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "userId", Value: mockUser.ID.String()}},
		)
		c.Set("user", &mockUserToken)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(map[string]any{
			"data": map[string]any{
				"name":         *updateInput.Name,
				"username":     *updateInput.Username,
				"email":        mockUser.Email,
				"country_code": *updateInput.CountryCode,
				"country":      "Brazil",
				"user_role":    mockUser.UserRole,
				"role":         *updateInput.Role,
				"interest":     *updateInput.Interest,
				"institution":  *updateInput.Institution,
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Missing user in context", func(t *testing.T) {
		svc := &mocks.MockUserService{}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/users/me",
			testutils.ToJSON(updateInput),
			nil,
			nil,
		)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateUserTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockUserService{}
			handler := user.NewUserHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPut,
				"/api/users/me",
				tt.Body,
				nil,
				nil,
			)
			c.Set("user", &mockUserToken)
			handler.UpdateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Username already exists", func(t *testing.T) {
		svc := &mocks.MockUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error) {
				return &models.UserResponse{}, nil
			},
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.UserUpdateInput,
				language string,
			) (*models.UserResponse, error) {
				return nil, services.ErrConflictUsername
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/users/me",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "userId", Value: mockUser.ID.String()}},
		)
		c.Set("user", &mockUserToken)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Username already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Country Code", func(t *testing.T) {
		svc := &mocks.MockUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error) {
				return &models.UserResponse{}, nil
			},
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.UserUpdateInput,
				language string,
			) (*models.UserResponse, error) {
				return nil, services.ErrInvalidCountryCode
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/users/me",
			testutils.ToJSON(updateInput),
			nil,
			nil,
		)
		c.Set("user", &mockUserToken)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "No country was found with this code.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		svc := &mocks.MockUserService{
			UpdateFunc: func(
				ctx context.Context,
				userID uuid.UUID,
				input models.UserUpdateInput,
				language string,
			) (*models.UserResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/users/me",
			testutils.ToJSON(updateInput),
			nil,
			nil,
		)
		c.Set("user", &mockUserToken)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(
			map[string]string{"error": "User not found."},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockUserService{
			UpdateFunc: func(
				ctx context.Context,
				userID uuid.UUID,
				input models.UserUpdateInput,
				language string,
			) (*models.UserResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := user.NewUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/users/me",
			testutils.ToJSON(updateInput),
			nil,
			nil,
		)
		c.Set("user", &mockUserToken)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
