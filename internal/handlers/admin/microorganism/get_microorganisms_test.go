package microorganism_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/microorganism"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetMicroorganisms(t *testing.T) {
	testutils.SetupTestContext()

	mockMicro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Escherichia coli",
		map[string]string{"pt": "Padrão", "en": "Standard", "es": "Estándar"},
		true,
	)

	mockResponse := mockMicro.ToAdminTableResponse("pt")

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindAllFunc: func(ctx context.Context, lang string) ([]models.MicroorganismAdminTableResponse, error) {
				return []models.MicroorganismAdminTableResponse{mockResponse}, nil
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/microorganism",
			"",
			nil,
			nil,
		)

		handler.GetMicroorganisms(c)

		expected := testutils.ToJSON(
			map[string][]models.MicroorganismAdminTableResponse{
				"data": {mockResponse},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindAllFunc: func(ctx context.Context, lang string) ([]models.MicroorganismAdminTableResponse, error) {
				return nil, gorm.ErrInvalidTransaction // Simula erro interno
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/microorganism",
			"",
			nil,
			nil,
		)

		handler.GetMicroorganisms(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}