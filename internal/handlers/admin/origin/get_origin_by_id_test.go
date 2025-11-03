package origin_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetOriginByID(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	id := uuid.NewString()
	mockOrigin := testmodels.NewOrigin(
		id,
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)
	db.Create(&mockOrigin)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: id}},
		)

		origin.GetOriginByID(c)

		expected := testutils.ToJSON(
			map[string]models.Origin{
				"data": mockOrigin,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: "132"}},
		)

		origin.GetOriginByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Origin not found", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)

		origin.GetOriginByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Origin not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("DB error", func(t *testing.T) {
		origRepo := repository.OriginRepo
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repository.OriginRepo = repository.NewOriginRepo(mockDB)
		defer func() {
			repository.OriginRepo = origRepo
		}()

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: id}},
		)

		origin.GetOriginByID(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
