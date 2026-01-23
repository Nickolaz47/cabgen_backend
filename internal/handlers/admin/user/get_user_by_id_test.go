package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	testutils.SetupTestContext()
	lang := "en"

	mockLoginUser := testmodels.NewLoginUser()
	response := mockLoginUser.ToAdminResponse(lang)
	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.AdminUserResponse, error) {
				return &response, nil
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/users", "",
			nil, gin.Params{{Key: "userId", Value: mockLoginUser.ID.String()}},
		)
		handler.GetUserByID(c)

		expected := testutils.ToJSON(
			map[string]models.AdminUserResponse{
				"data": response,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/users", "",
			nil, gin.Params{{Key: "userId", Value: "ifew90843"}},
		)
		handler.GetUserByID(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "The URL ID is invalid.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.AdminUserResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/users", "",
			nil, gin.Params{{Key: "userId", Value: uuid.NewString()}},
		)
		handler.GetUserByID(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "User not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID, language string) (*models.AdminUserResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/users", "",
			nil, gin.Params{{Key: "userId", Value: uuid.NewString()}},
		)
		handler.GetUserByID(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
