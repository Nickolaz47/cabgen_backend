package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetUsernames(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindAllUsernamesFunc: func(ctx context.Context) ([]string, error) {
				return []string{"alice", "nick"}, nil
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users/usernames",
			"", nil, nil,
		)
		handler.GetUsernames(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []string{"alice", "nick"},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockAdminUserService{
			FindAllUsernamesFunc: func(ctx context.Context) ([]string, error) {
				return nil, services.ErrInternal
			},
		}
		handler := user.NewAdminUserHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/users/usernames",
			"", nil, nil,
		)
		handler.GetUsernames(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
