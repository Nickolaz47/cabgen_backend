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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testutils.SetupTestContext()

	mockAdminUser := testmodels.NewAdminLoginUser()
	mockAdminUserToken := testmodels.NewUserToken(
		mockAdminUser.ID,
		mockAdminUser.Username,
		mockAdminUser.UserRole,
	)

	createInput := testmodels.NewAdminCreateUserInput("", "")

	createResponse := models.AdminUserResponse{
		ID:          uuid.New(),
		Name:        createInput.Name,
		Username:    createInput.Username,
		Email:       createInput.Email,
		CountryCode: createInput.CountryCode,
		Country:     "Brazil",
		UserRole:    models.Collaborator,
		CreatedBy:   mockAdminUser.Username,
		ActivatedBy: &mockAdminUser.Username,
		ActivatedOn: &time.Time{},
		IsActive:    true,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			CreateFunc: func(
				ctx context.Context,
				input models.AdminUserCreateInput,
				adminName, language string,
			) (*models.AdminUserResponse, error) {
				return &createResponse, nil
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/users",
			testutils.ToJSON(createInput),
			nil,
			nil,
		)
		c.Set("user", &mockAdminUserToken)
		handler.CreateUser(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "User created successfully.",
			"data": map[string]any{
				"name":         createInput.Name,
				"username":     createInput.Username,
				"email":        createInput.Email,
				"country_code": createInput.CountryCode,
				"country":      "Brazil",
				"is_active":    true,
				"user_role":    "Collaborator",
				"updated_at":   time.Time{},
				"created_at":   time.Time{},
				"activated_on": time.Time{},
				"activated_by": mockAdminUser.Username,
				"created_by":   mockAdminUser.Username,
			},
		})

		var got map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)

		if data, ok := got["data"].(map[string]any); ok {
			delete(data, "id")
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, testutils.ToJSON(got))
	})

	t.Run("Error - Missing user in context", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/users",
			testutils.ToJSON(createInput),
			nil,
			nil,
		)
		handler.CreateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue."},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.AdminCreateUserTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockAdminUserService{}
			handler := user.NewAdminUserHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPost,
				"/api/admin/users",
				tt.Body,
				nil,
				nil,
			)
			c.Set("user", &mockAdminUserToken)
			handler.CreateUser(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Email already exists", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			CreateFunc: func(
				ctx context.Context,
				input models.AdminUserCreateInput,
				adminName, language string,
			) (*models.AdminUserResponse, error) {
				return nil, services.ErrConflictEmail
			},
		}
		handler := user.NewAdminUserHandler(svc)

		body := testutils.ToJSON(
			testmodels.NewAdminCreateUserInput("cadmin@mail.com", ""),
		)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/user",
			body,
			nil,
			nil,
		)
		c.Set("user", &mockAdminUserToken)
		handler.CreateUser(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Email is already in use.",
		})

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Username already exists", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			CreateFunc: func(
				ctx context.Context,
				input models.AdminUserCreateInput,
				adminName, language string,
			) (*models.AdminUserResponse, error) {
				return nil, services.ErrConflictUsername
			},
		}
		handler := user.NewAdminUserHandler(svc)

		body := testutils.ToJSON(
			testmodels.NewAdminCreateUserInput("", "admin"),
		)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/user",
			body,
			nil,
			nil,
		)
		c.Set("user", &mockAdminUserToken)

		handler.CreateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Username already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			CreateFunc: func(
				ctx context.Context,
				input models.AdminUserCreateInput,
				adminName, language string,
			) (*models.AdminUserResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := user.NewAdminUserHandler(svc)

		body := testutils.ToJSON(
			testmodels.NewAdminCreateUserInput("", "admin"),
		)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/user",
			body,
			nil,
			nil,
		)
		c.Set("user", &mockAdminUserToken)
		handler.CreateUser(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again."},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
