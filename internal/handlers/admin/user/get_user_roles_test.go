package user_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/user"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetUserRoles(t *testing.T) {
	testutils.SetupTestContext()

	svc := &mocks.MockAdminUserService{}
	handler := user.NewAdminUserHandler(svc)

	c, w := testutils.SetupGinContext(
		http.MethodGet, "/api/admin/users/roles",
		"", nil, nil,
	)
	handler.GetUserRoles(c)

	expected := testutils.ToJSON(
		map[string]any{
			"data": models.UserRoles,
		},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
