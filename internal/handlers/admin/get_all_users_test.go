package admin_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockCountry := testmodels.NewCountry("", "", "", "")
	db.Create(&mockCountry)

	mockLoginUser := testmodels.NewLoginUser()
	db.Create(&mockLoginUser)

	c, w := testutils.SetupGinContext(
		http.MethodGet, "/api/admin/user", "",
		nil, nil,
	)

	admin.GetAllUsers(c)

	expected := testutils.ToJSON(
		map[string][]models.User{
			"data": {mockLoginUser},
		},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
