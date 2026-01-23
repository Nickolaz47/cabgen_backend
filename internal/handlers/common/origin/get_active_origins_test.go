package origin

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetActiveOrigins(t *testing.T) {
	testutils.SetupTestContext()
	mockOrigin := models.Origin{
		ID:       uuid.New(),
		Names:    map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockOriginService{
			FindAllActiveFunc: func(ctx context.Context, lang string) ([]models.OriginFormResponse, error) {
				return []models.OriginFormResponse{mockOrigin.ToFormResponse(lang)}, nil
			},
		}
		handler := NewOriginHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/origin", "",
			nil, nil,
		)
		handler.GetActiveOrigins(c)

		expected := testutils.ToJSON(map[string][]models.OriginFormResponse{
			"data": {mockOrigin.ToFormResponse("en")},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockOriginService{
			FindAllActiveFunc: func(ctx context.Context, lang string) ([]models.OriginFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := NewOriginHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/origin", "",
			nil, nil,
		)
		handler.GetActiveOrigins(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
