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
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateSampleSource(t *testing.T) {
	testutils.SetupTestContext()

	mockSampleSourceInput := models.SampleSourceCreateInput{
		Names:    map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		Groups:   map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			CreateFunc: func(ctx context.Context, input models.SampleSourceCreateInput) (*models.SampleSourceAdminDetailResponse, error) {
				return &models.SampleSourceAdminDetailResponse{
					Names:    mockSampleSourceInput.Names,
					Groups:   mockSampleSourceInput.Groups,
					IsActive: mockSampleSourceInput.IsActive,
				}, nil
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sample-source", body,
			nil, nil,
		)
		handler.CreateSampleSource(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "Sample source created successfully.",
			"data": models.SampleSourceAdminDetailResponse{
				Names:    mockSampleSourceInput.Names,
				Groups:   mockSampleSourceInput.Groups,
				IsActive: mockSampleSourceInput.IsActive,
			},
		})

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.CreateSampleSourceTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockSampleSourceService{}
			handler := samplesource.NewAdminSampleSourceHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/sample-source", tt.Body,
				nil, nil,
			)
			handler.CreateSampleSource(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			CreateFunc: func(ctx context.Context, input models.SampleSourceCreateInput) (*models.SampleSourceAdminDetailResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/origin", testutils.ToJSON(mockSampleSourceInput),
			nil, nil,
		)
		handler.CreateSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Sample source already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			CreateFunc: func(ctx context.Context, input models.SampleSourceCreateInput) (*models.SampleSourceAdminDetailResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sample-source", body,
			nil, nil,
		)
		handler.CreateSampleSource(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
