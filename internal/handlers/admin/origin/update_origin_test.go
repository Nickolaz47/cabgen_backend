package origin_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpdateOrigin(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	id := uuid.NewString()
	mockOrigin := testmodels.NewOrigin(
		id,
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)
	db.Create(&mockOrigin)

	isActive := false
	mockOriginInput := models.OriginUpdateInput{
		IsActive: &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		body := testutils.ToJSON(mockOriginInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/origin", body,
			nil, gin.Params{{Key: "originId", Value: id}},
		)

		origin.UpdateOrigin(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": map[string]any{
					"id":        id,
					"name":      mockOrigin.Names["en"],
					"is_active": *mockOriginInput.IsActive,
				},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateOriginTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/origin", tt.Body,
				nil, gin.Params{{Key: "originId", Value: id}},
			)

			origin.UpdateOrigin(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: "12"}},
		)

		origin.UpdateOrigin(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Origin not found", func(t *testing.T) {
		body := testutils.ToJSON(mockOriginInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/origin", body,
			nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)

		origin.UpdateOrigin(c)

		expected := testutils.ToJSON(
			map[string]string{"error": "Origin not found."},
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

		body := testutils.ToJSON(mockOriginInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/origin", body,
			nil, gin.Params{{Key: "originId", Value: id}},
		)

		origin.UpdateOrigin(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
