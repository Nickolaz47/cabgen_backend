package analysis_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestGetAnalysisByID(t *testing.T) {
	testutils.SetupTestContext()

	mockAnalysis := testmodels.CreateMockAnalysis()
	mockResponse := mockAnalysis.ToResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockResponse,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: "abc1"}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.User.ID.String()}},
		)
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not found", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.User.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Analysis not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.User.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}

func TestGetAnalysisFastQCByID(t *testing.T) {
	testutils.SetupTestContext()

	mockAnalysis := testmodels.CreateMockAnalysis()

	f1, err := os.CreateTemp("", "fastqc1-*.html")
	assert.NoError(t, err)
	defer os.Remove(f1.Name())
	f1.WriteString("<html>fastqc1</html>")
	f1.Close()

	f2, err := os.CreateTemp("", "fastqc2-*.html")
	assert.NoError(t, err)
	defer os.Remove(f2.Name())
	f2.WriteString("<html>fastqc2</html>")
	f2.Close()

	mockAnalysis.FastQC1 = strPtr(f1.Name())
	mockAnalysis.FastQC2 = strPtr(f2.Name())
	mockResponse := mockAnalysis.ToResponse()

	t.Run("Success - fastqc1", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return &mockResponse, nil
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis/fastqc1", "", nil,
			gin.Params{
				{Key: "analysisId", Value: mockAnalysis.ID.String()},
				{Key: "fastqcReport", Value: "fastqc1"},
			},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "fastqc1")
	})

	t.Run("Success - fastqc2", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return &mockResponse, nil
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis/fastqc2", "", nil,
			gin.Params{
				{Key: "analysisId", Value: mockAnalysis.ID.String()},
				{Key: "fastqcReport", Value: "fastqc2"},
			},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "fastqc2")
	})

	t.Run("Error - Invalid UUID", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis/fastqc1", "", nil,
			gin.Params{
				{Key: "analysisId", Value: "abc1"},
				{Key: "fastqcReport", Value: "fastqc1"},
			},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - No user token", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis/fastqc1", "", nil,
			gin.Params{
				{Key: "analysisId", Value: mockAnalysis.ID.String()},
				{Key: "fastqcReport", Value: "fastqc1"},
			},
		)
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Not found", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis/fastqc1", "", nil,
			gin.Params{
				{Key: "analysisId", Value: mockAnalysis.ID.String()},
				{Key: "fastqcReport", Value: "fastqc1"},
			},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error - Invalid fastqcReport param", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return &mockResponse, nil
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis/invalid", "", nil,
			gin.Params{
				{Key: "analysisId", Value: mockAnalysis.ID.String()},
				{Key: "fastqcReport", Value: "invalid"},
			},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error - FastQC report not available", func(t *testing.T) {
		nilResponse := mockAnalysis.ToResponse()
		nilResponse.FastQC1 = nil
		nilResponse.FastQC2 = nil

		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return &nilResponse, nil
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis/fastqc1", "", nil,
			gin.Params{
				{Key: "analysisId", Value: mockAnalysis.ID.String()},
				{Key: "fastqcReport", Value: "fastqc1"},
			},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis/fastqc1", "", nil,
			gin.Params{
				{Key: "analysisId", Value: mockAnalysis.ID.String()},
				{Key: "fastqcReport", Value: "fastqc1"},
			},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
