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

func TestGetMicroorganismBySpecies(t *testing.T) {
	testutils.SetupTestContext()

	mockMicro1 := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Salmonella",
		map[string]string{"pt": "ssp", "en": "ssp", "es": "ssp"},
		true,
	)

	mockMicro2 := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Escherichia coli",
		map[string]string{"pt": "O157:H7", "en": "O157:H7", "es": "O157:H7"},
		true,
	)

	language := "en"
	mockResponse1 := mockMicro1.ToAdminTableResponse(language)
	mockResponse2 := mockMicro2.ToAdminTableResponse(language)

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindBySpeciesFunc: func(ctx context.Context, species, lang string) ([]models.MicroorganismAdminTableResponse, error) {
				return []models.MicroorganismAdminTableResponse{mockResponse1}, nil
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		species := "salmonella"
		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/microorganism/search?species="+species,
			"",
			nil,
			nil,
		)

		handler.GetMicroorganismBySpecies(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.MicroorganismAdminTableResponse{mockResponse1},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Input Empty", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindAllFunc: func(ctx context.Context, lang string) ([]models.MicroorganismAdminTableResponse, error) {
				return []models.MicroorganismAdminTableResponse{mockResponse1, mockResponse2}, nil
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/microorganism/search?species=",
			"",
			nil,
			nil,
		)

		handler.GetMicroorganismBySpecies(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.MicroorganismAdminTableResponse{mockResponse1, mockResponse2},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindBySpeciesFunc: func(ctx context.Context, species, lang string) ([]models.MicroorganismAdminTableResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		species := "error"
		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/microorganism/search?species="+species,
			"",
			nil,
			nil,
		)

		handler.GetMicroorganismBySpecies(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
