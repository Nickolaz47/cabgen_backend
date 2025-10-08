package origin_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateOrigin(t *testing.T) {
	testutils.SetupTestContext()
	testutils.SetupTestRepos()

	mockOriginInput := models.OriginCreateInput{
		Names:    map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		body := testutils.ToJSON(mockOriginInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/origin", body,
			nil, nil,
		)

		origin.CreateOrigin(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "Origin created successfully.",
			"data": map[string]any{
				"name":      mockOriginInput.Names["en"],
				"is_active": mockOriginInput.IsActive,
			},
		})

		var result map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)

		if data, ok := result["data"].(map[string]any); ok {
			delete(data, "id")
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, testutils.ToJSON(result))
	})

	for _, tt := range data.CreateOriginTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/origin", tt.Body,
				nil, nil,
			)

			origin.CreateOrigin(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

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
			http.MethodPost, "/api/admin/origin", body,
			nil, nil,
		)

		origin.CreateOrigin(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
