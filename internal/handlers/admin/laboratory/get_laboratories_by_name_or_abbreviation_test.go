package laboratory_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetLaboratoriesByNameOrAbbreviation(t *testing.T) {
	testutils.SetupTestContext()

	lab1 := testmodels.NewLaboratory(
		uuid.NewString(),
		"Laboratório Dom Bosco",
		"LABDB",
		true,
	)

	lab2 := testmodels.NewLaboratory(
		uuid.NewString(),
		"Laboratório Bittar",
		"LABB",
		true,
	)

	adminResponse1 := lab1.ToAdminTableResponse()
	adminResponse2 := lab2.ToAdminTableResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockLaboratoryService{
			FindByNameOrAbbreviationFunc: func(
				ctx context.Context,
				input string,
			) ([]models.LaboratoryAdminTableResponse, error) {
				return []models.LaboratoryAdminTableResponse{
					adminResponse1,
				}, nil
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/laboratory/search?nameOrAbbreaviation=dom",
			"",
			nil,
			nil,
		)
		handler.GetLaboratoriesByNameOrAbbreviation(c)

		expected := testutils.ToJSON(
			map[string][]models.LaboratoryAdminTableResponse{
				"data": {adminResponse1},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Input Empty", func(t *testing.T) {
		svc := &mocks.MockLaboratoryService{
			FindAllFunc: func(
				ctx context.Context,
			) ([]models.LaboratoryAdminTableResponse, error) {
				return []models.LaboratoryAdminTableResponse{
					adminResponse1,
					adminResponse2,
				}, nil
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/laboratory/search?nameOrAbbreaviation=",
			"",
			nil,
			nil,
		)
		handler.GetLaboratoriesByNameOrAbbreviation(c)

		expected := testutils.ToJSON(
			map[string][]models.LaboratoryAdminTableResponse{
				"data": {adminResponse1, adminResponse2},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockLaboratoryService{
			FindByNameOrAbbreviationFunc: func(
				ctx context.Context,
				input string,
			) ([]models.LaboratoryAdminTableResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/laboratory/search?nameOrAbbreaviation=dom",
			"",
			nil,
			nil,
		)
		handler.GetLaboratoriesByNameOrAbbreviation(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
