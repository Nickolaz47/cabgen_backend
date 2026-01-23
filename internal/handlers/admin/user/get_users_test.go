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
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	testutils.SetupTestContext()
	lang := "en"

	mockUser := testmodels.NewLoginUser()
	mockUser2 := testmodels.NewAdminLoginUser()

	userResponse := mockUser.ToAdminResponse(lang)
	userResponse2 := mockUser2.ToAdminResponse(lang)

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindFunc: func(ctx context.Context, filter models.AdminUserFilter, language string) ([]models.AdminUserResponse, error) {
				return []models.AdminUserResponse{userResponse, userResponse2}, nil
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users",
			"", nil, nil,
		)
		handler.GetUsers(c)

		expected := testutils.ToJSON(
			map[string][]models.AdminUserResponse{
				"data": {userResponse, userResponse2},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - With filter", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindFunc: func(ctx context.Context, filter models.AdminUserFilter, language string) ([]models.AdminUserResponse, error) {
				return []models.AdminUserResponse{userResponse2}, nil
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users",
			"", nil, gin.Params{
				{Key: "userRole", Value: "admin"},
			},
		)
		handler.GetUsers(c)

		expected := testutils.ToJSON(
			map[string][]models.AdminUserResponse{
				"data": {userResponse2},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindFunc: func(ctx context.Context, filter models.AdminUserFilter, language string) ([]models.AdminUserResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users",
			"", nil, gin.Params{
				{Key: "userRole", Value: "admin"},
			},
		)
		handler.GetUsers(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
