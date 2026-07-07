package selectoptions_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/selectoptions"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetSelectOptions(t *testing.T) {
	testutils.SetupTestContext()

	mockResponse := models.EnumSelectsResponse{
		Roles: []models.SelectOption{
			{Label: "option.role.admin", Value: "Admin"},
			{Label: "option.role.collaborator", Value: "Collaborator"},
		},
		Taxons: []models.SelectOption{
			{Label: "option.taxon.bacteria", Value: "Bacteria"},
			{Label: "option.taxon.fungi", Value: "Fungi"},
			{Label: "option.taxon.protozoa", Value: "Protozoa"},
			{Label: "option.taxon.virus", Value: "Virus"},
		},
		Genders: []models.SelectOption{
			{Label: "option.gender.female", Value: "Female"},
			{Label: "option.gender.male", Value: "Male"},
			{Label: "option.gender.unspecified", Value: "Unspecified"},
		},
		HealthServiceTypes: []models.SelectOption{
			{Label: "option.health_service_type.public", Value: "Public"},
			{Label: "option.health_service_type.private", Value: "Private"},
		},
		AnalysisTypes: []models.SelectOption{
			{Label: "option.analysis_type.fastqc", Value: "FASTQC"},
			{Label: "option.analysis_type.genome", Value: "GENOME"},
			{Label: "option.analysis_type.complete", Value: "COMPLETE"},
		},
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSelectOptionsService{
			FindAllFunc: func(ctx context.Context) (
				*models.EnumSelectsResponse, error) {
				return &mockResponse, nil
			},
		}
		handler := selectoptions.NewSelectOptionsHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/select-options",
			"", nil, nil,
		)
		handler.GetSelectOptions(c)

		expected := testutils.ToJSON(map[string]any{
			"data": mockResponse,
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockSelectOptionsService{
			FindAllFunc: func(ctx context.Context) (*models.EnumSelectsResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := selectoptions.NewSelectOptionsHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/select-options",
			"", nil, nil,
		)
		handler.GetSelectOptions(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
