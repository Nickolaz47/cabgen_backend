package samplesource_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetSampleSourceByID(t *testing.T) {
	testutils.SetupTestContext()

	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSourceAdminDetailResponse, error) {
				response := mockSampleSource.ToAdminDetailResponse()
				return &response, nil
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sample-source", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
		)
		handler.GetSampleSourceByID(c)

		expected := testutils.ToJSON(
			map[string]models.SampleSourceAdminDetailResponse{
				"data": mockSampleSource.ToAdminDetailResponse(),
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sample-source", "",
			nil, nil,
		)
		handler.GetSampleSourceByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Sample source not found", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSourceAdminDetailResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sample-source", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)
		handler.GetSampleSourceByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Sample source not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Database error", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SampleSourceAdminDetailResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sample-source", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)
		handler.GetSampleSourceByID(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
