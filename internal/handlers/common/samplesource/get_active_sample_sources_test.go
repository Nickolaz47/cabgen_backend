package samplesource_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetActiveSampleSources(t *testing.T) {
	testutils.SetupTestContext()
	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			FindAllActiveFunc: func(ctx context.Context, language string) ([]models.SampleSourceFormResponse, error) {
				return []models.SampleSourceFormResponse{mockSampleSource.ToFormResponse("en")}, nil
			},
		}
		handler := samplesource.NewSampleSourceHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/sample-source",
			"", nil, nil,
		)
		handler.GetActiveSampleSources(c)

		expected := testutils.ToJSON(map[string][]models.SampleSourceFormResponse{
			"data": {mockSampleSource.ToFormResponse("en")},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			FindAllActiveFunc: func(ctx context.Context, language string) ([]models.SampleSourceFormResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := samplesource.NewSampleSourceHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/sample-source",
			"", nil, nil,
		)
		handler.GetActiveSampleSources(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
