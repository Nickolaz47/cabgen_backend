package analysis_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/analysis"
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
	mockResponse := mockAnalysis.ToResponse()

	mockUserID := uuid.New()

	validInput := map[string]any{
		"type":      models.AnalysisTypeFastQC,
		"sample_id": validUUID,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			CreateFunc: func(ctx context.Context,
				input models.AnalysisCreateDTO) (
				*models.AnalysisResponse, error) {
				assert.Equal(t, mockUserID, input.UserID)
				return &mockResponse, nil
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analysis",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

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
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		invalidInput := testutils.CopyMap(validInput)
		invalidInput["type"] = "invalid_type"

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analysis",
			testutils.ToJSON(invalidInput),
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

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
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		for _, test := range data.AnalysisCreateTests {
			t.Run(test.Name, func(t *testing.T) {
				c, w := testutils.SetupGinContext(
					http.MethodPost,
					"/api/analysis",
					test.Body,
					nil,
					nil,
				)
				c.Set("user", &models.UserToken{ID: mockUserID})

				handler.CreateAnalysis(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, test.Expected, w.Body.String())
			})
		}
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analysis",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)

		handler.CreateAnalysis(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			CreateFunc: func(ctx context.Context,
				input models.AnalysisCreateDTO) (
				*models.AnalysisResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analysis",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

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
