package origin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetOriginsByName(t *testing.T) {
	testutils.SetupTestContext()

	mockOrigin := testmodels.NewOrigin(
		uuid.NewString(),
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)

	mockOrigin2 := testmodels.NewOrigin(
		uuid.NewString(),
		map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		true,
	)

	language := "en"
	mockResponse1 := mockOrigin.ToAdminTableResponse(language)
	mockResponse2 := mockOrigin2.ToAdminTableResponse(language)

	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockOriginService{
			FindByNameFunc: func(ctx context.Context, name, lang string) ([]models.OriginAdminTableResponse, error) {
				return []models.OriginAdminTableResponse{mockResponse1}, nil
			},
		}

		handler := origin.NewAdminOriginHandler(&svc)

		name := "food"
		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin/search?name="+name,
			"",
			nil,
			nil,
		)

		handler.GetOriginsByName(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.OriginAdminTableResponse{mockResponse1},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Input Empty", func(t *testing.T) {
		svc := testmodels.MockOriginService{
			FindAllFunc: func(ctx context.Context, lang string) ([]models.OriginAdminTableResponse, error) {
				return []models.OriginAdminTableResponse{mockResponse1, mockResponse2}, nil
			},
		}

		handler := origin.NewAdminOriginHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin/search?name=",
			"",
			nil,
			nil,
		)

		handler.GetOriginsByName(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.OriginAdminTableResponse{mockResponse1, mockResponse2},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := testmodels.MockOriginService{
			FindByNameFunc: func(ctx context.Context, name, lang string) ([]models.OriginAdminTableResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := origin.NewAdminOriginHandler(&svc)

		name := "human"
		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin/search?name="+name,
			"",
			nil,
			nil,
		)

		handler.GetOriginsByName(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
