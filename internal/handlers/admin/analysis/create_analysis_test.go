package analysis_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAnalysis(t *testing.T) {
	testutils.SetupTestContext()

	const validUUID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	mockAnalysis := testmodels.CreateMockAnalysis()
	mockResponse := mockAnalysis.ToAdminResponse()

	mockUserID := uuid.New()

	validInput := map[string]any{
		"type":      models.AnalysisTypeFastQC,
		"sample_id": validUUID,
		"user_id":   mockUserID,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			CreateFunc: func(ctx context.Context,
				input models.AnalysisCreateDTO) (
				*models.AnalysisAdminResponse, error) {
				assert.Equal(t, mockUserID, input.UserID)
				return &mockResponse, nil
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/analysis",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		handler.CreateAnalysis(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data":    mockResponse,
				"message": "Analysis created successfully.",
			},
		)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Type", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{}
		handler := analysis.NewAdminAnalysisHandler(svc)

		invalidInput := testutils.CopyMap(validInput)
		invalidInput["type"] = "invalid_type"

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/analysis",
			testutils.ToJSON(invalidInput),
			nil,
			nil,
		)
		handler.CreateAnalysis(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "This analysis type is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Bad Request", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{}
		handler := analysis.NewAdminAnalysisHandler(svc)

		for _, test := range data.AdminAnalysisCreateTests {
			t.Run(test.Name, func(t *testing.T) {
				c, w := testutils.SetupGinContext(
					http.MethodPost,
					"/api/admin/analysis",
					test.Body,
					nil,
					nil,
				)
				handler.CreateAnalysis(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, test.Expected, w.Body.String())
			})
		}
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			CreateFunc: func(ctx context.Context,
				input models.AnalysisCreateDTO) (
				*models.AnalysisAdminResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/analysis",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		handler.CreateAnalysis(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
