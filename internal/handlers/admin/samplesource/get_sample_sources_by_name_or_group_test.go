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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetSampleSourceByNameOrGroup(t *testing.T) {
	testutils.SetupTestContext()

	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)
	lang := "en"

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			FindByNameOrGroupFunc: func(ctx context.Context, input, language string) ([]models.SampleSourceAdminTableResponse, error) {
				return []models.SampleSourceAdminTableResponse{
					mockSampleSource.ToAdminTableResponse(lang),
				}, nil
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sample-source/search?nameOrGroup=plas", "",
			nil, nil,
		)
		handler.GetSampleSourcesByNameOrGroup(c)

		expected := testutils.ToJSON(map[string]any{
			"data": []models.SampleSourceAdminTableResponse{
				mockSampleSource.ToAdminTableResponse(lang),
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Input empty", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			FindAllFunc: func(ctx context.Context, language string) ([]models.SampleSourceAdminTableResponse, error) {
				return []models.SampleSourceAdminTableResponse{
					mockSampleSource.ToAdminTableResponse(lang),
				}, nil
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sample-source/search?nameOrGroup=", "",
			nil, nil,
		)
		handler.GetSampleSourcesByNameOrGroup(c)

		expected := testutils.ToJSON(map[string][]models.SampleSourceAdminTableResponse{
			"data": []models.SampleSourceAdminTableResponse{
				mockSampleSource.ToAdminTableResponse(lang),
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Database error", func(t *testing.T) {
		svc := &mocks.MockSampleSourceService{
			FindByNameOrGroupFunc: func(ctx context.Context, input, language string) ([]models.SampleSourceAdminTableResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sample-source/search?nameOrGroup=blo", "",
			nil, nil,
		)
		handler.GetSampleSourcesByNameOrGroup(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
