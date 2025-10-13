package origin_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetActiveOrigins(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockOrigin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	mockOrigin2 := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"}, false)
	db.Create(&mockOrigin)
	db.Create(&mockOrigin2)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/origin", "",
			nil, nil,
		)

		origin.GetActiveOrigins(c)

		expected := testutils.ToJSON(map[string]any{
			"data": []models.OriginPublicResponse{
				mockOrigin.ToPublicResponse(c),
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		origRepo := repository.OriginRepo
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repository.OriginRepo = repository.NewOriginRepo(mockDB)
		defer func() {
			repository.OriginRepo = origRepo
		}()

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/origin", "",
			nil, nil,
		)

		origin.GetActiveOrigins(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
