package user_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
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

	mockAdminUser := testmodels.NewAdminLoginUser()
	updateInput := testmodels.NewAdminUpdateUserInput()

	updateResponse := models.AdminUserResponse{
		ID:          uuid.New(),
		Name:        *updateInput.Name,
		Username:    *updateInput.Username,
		Email:       *updateInput.Email,
		CountryCode: *updateInput.CountryCode,
		Country:     "Brazil",
		UserRole:    *updateInput.UserRole,
		Role:        updateInput.Role,
		Interest:    updateInput.Interest,
		Institution: updateInput.Institution,
		ActivatedBy: &mockAdminUser.Username,
		ActivatedOn: &time.Time{},
		CreatedBy:   mockAdminUser.Username,
		IsActive:    true,
	}

	validUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.AdminUserUpdateInput,
				language string,
			) (*models.AdminUserResponse, error) {
				return &updateResponse, nil
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/users",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "userId", Value: validUserID.String()}},
		)
		handler.UpdateUser(c)

		var got map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)

		if data, ok := got["data"].(map[string]any); ok {
			delete(data, "id")
		}

		expected := testutils.ToJSON(map[string]any{
			"data": map[string]any{
				"name":         *updateInput.Name,
				"username":     *updateInput.Username,
				"email":        *updateInput.Email,
				"country_code": *updateInput.CountryCode,
				"country":      "Brazil",
				"user_role":    updateInput.UserRole,
				"role":         *updateInput.Role,
				"interest":     *updateInput.Interest,
				"institution":  *updateInput.Institution,
				"created_at":   time.Time{},
				"activated_by": mockAdminUser.Username,
				"created_by":   mockAdminUser.Username,
				"activated_on": time.Time{},
				"updated_at":   time.Time{},
				"is_active":    true,
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, testutils.ToJSON(got))
	})

	t.Run("Error - Invalid user ID", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/users",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "userId", Value: "invalid-id"}},
		)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(
			map[string]string{"error": "The URL ID is invalid."})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.AdminUpdateUserTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockAdminUserService{}
			handler := user.NewAdminUserHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPut,
				"/api/admin/users",
				tt.Body,
				nil,
				gin.Params{{Key: "userId", Value: validUserID.String()}},
			)

			handler.UpdateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - User not found", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.AdminUserUpdateInput,
				language string,
			) (*models.AdminUserResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/users",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "userId", Value: validUserID.String()}},
		)
		handler.UpdateUser(c)

		expected := `{"error":"User not found."}`

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Username already exists", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.AdminUserResponse, error) {
				return &models.AdminUserResponse{}, nil
			},
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.AdminUserUpdateInput,
				language string,
			) (*models.AdminUserResponse, error) {
				return nil, services.ErrConflictUsername
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/users",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "userId", Value: validUserID.String()}},
		)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Username already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Email already exists", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.AdminUserResponse, error) {
				return &models.AdminUserResponse{}, nil
			},
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.AdminUserUpdateInput,
				language string,
			) (*models.AdminUserResponse, error) {
				return nil, services.ErrConflictEmail
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/users",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "userId", Value: validUserID.String()}},
		)
		handler.UpdateUser(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Email is already in use.",
		})

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.AdminUserUpdateInput,
				language string,
			) (*models.AdminUserResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/user/"+validUserID.String(),
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "userId", Value: validUserID.String()}},
		)
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
