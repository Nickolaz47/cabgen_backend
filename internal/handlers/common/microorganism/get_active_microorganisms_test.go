package microorganism

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

func TestGetActiveMicroorganisms(t *testing.T) {
	testutils.SetupTestContext()
	mockMicro := models.Microorganism{
		ID:       uuid.New(),
		Species:  "E. coli",
		Variety:  map[string]string{"pt": "Variedade A", "en": "Variety A"},
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindAllActiveFunc: func(ctx context.Context, lang string) ([]models.MicroorganismFormResponse, error) {
				return []models.MicroorganismFormResponse{mockMicro.ToFormResponse(lang)}, nil
			},
		}
		handler := NewMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/microorganisms", "",
			nil, nil,
		)
		handler.GetActiveMicroorganisms(c)

		expected := testutils.ToJSON(map[string][]models.MicroorganismFormResponse{
			"data": {mockMicro.ToFormResponse("en")},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindAllActiveFunc: func(ctx context.Context, lang string) ([]models.MicroorganismFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := NewMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/microorganisms", "",
			nil, nil,
		)
		handler.GetActiveMicroorganisms(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
