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

func TestGetOriginByName(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockOrigin := testmodels.NewOrigin(
		uuid.NewString(),
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)
	db.Create(&mockOrigin)

	t.Run("Success", func(t *testing.T) {
		name := "food"
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin/search?originName="+name, "",
			nil, nil,
		)

		origin.GetOriginByName(c)

		expected := testutils.ToJSON(
			map[string]models.Origin{
				"data": mockOrigin,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - originName empty", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin/search", "",
			nil, nil,
		)

		origin.GetOriginByName(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The search parameter originName is empty.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Origin not found", func(t *testing.T) {
		name := "human"
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin/search?originName="+name, "",
			nil, nil,
		)

		origin.GetOriginByName(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Origin not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
