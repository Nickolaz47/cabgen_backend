package origin_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllOrigins(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockOrigin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)
	db.Create(&mockOrigin)

	c, w := testutils.SetupGinContext(
		http.MethodGet, "/api/admin/origin", "",
		nil, nil,
	)

	origin.GetAllOrigins(c)

	expected := testutils.ToJSON(
		map[string][]models.Origin{
			"data": {mockOrigin},
		},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
