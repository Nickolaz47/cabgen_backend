package samplesource_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSampleSource(t *testing.T) {
	testutils.SetupTestContext()

	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasm", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sange"},
		false,
	)

	isActive := true
	mockSampleSourceInput := models.SampleSourceUpdateInput{
		Names:    map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		Groups:   map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		IsActive: &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SampleSourceUpdateInput) (*models.SampleSourceAdminDetailResponse, error) {
				return &models.SampleSourceAdminDetailResponse{
					ID:       mockSampleSource.ID,
					Names:    mockSampleSourceInput.Names,
					Groups:   mockSampleSourceInput.Groups,
					IsActive: isActive,
				}, nil
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sample-source", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
		)
		handler.UpdateSampleSource(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": models.SampleSourceAdminDetailResponse{
					ID:       mockSampleSource.ID,
					Names:    mockSampleSourceInput.Names,
					Groups:   mockSampleSourceInput.Groups,
					IsActive: isActive,
				},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateSampleSourceTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := testmodels.MockSampleSourceService{}
			handler := samplesource.NewAdminSampleSourceHandler(&svc)

			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/sample-source", tt.Body,
				nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
			)
			handler.UpdateSampleSource(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sample-source", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: "123"}},
		)
		handler.UpdateSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not found", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SampleSourceUpdateInput) (*models.SampleSourceAdminDetailResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sample-source", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)
		handler.UpdateSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Sample source not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SampleSourceUpdateInput) (*models.SampleSourceAdminDetailResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sample-source", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)
		handler.UpdateSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Sample source already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SampleSourceUpdateInput) (*models.SampleSourceAdminDetailResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sample-source", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
		)
		handler.UpdateSampleSource(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
