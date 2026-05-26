package sample_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateSample(t *testing.T) {
	testutils.SetupTestContext()

	const validUUID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	mockSample := testmodels.CreateMockSample()
	mockResponse := mockSample.ToResponse("")

	validInput := map[string]any{
		"name":              "Sample-SARS-CoV-2",
		"collection_date":   "2026-05-20T00:00:00Z",
		"run_number":        "RUN-2026-XYZ",
		"run_date":          "2026-05-25T00:00:00Z",
		"city":              "Maricá",
		"origin_code":       "BR-RJ-01",
		"gender":            "Male",
		"date_of_birth":     "1990-01-01T00:00:00Z",
		"country_code":      "BRA",
		"user_id":           validUUID,
		"origin_id":         validUUID,
		"sample_source_id":  validUUID,
		"microorganism_id":  validUUID,
		"sequencer_id":      validUUID,
		"laboratory_id":     validUUID,
		"health_service_id": validUUID,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			CreateFunc: func(ctx context.Context,
				input models.SampleCreateInput,
				language string) (*models.SampleResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/sample",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)

		handler.CreateSample(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data":    mockResponse,
				"message": "Sample created successfully.",
			},
		)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Minimal Payload", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			CreateFunc: func(ctx context.Context,
				input models.SampleCreateInput, language string) (
				*models.SampleResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := sample.NewAdminSampleHandler(svc)

		minimalInput := map[string]any{
			"name":              "Minimal-Sample",
			"collection_date":   "2026-05-20T00:00:00Z",
			"run_number":        "RUN-01",
			"run_date":          "2026-05-25T00:00:00Z",
			"country_code":      "BRA",
			"user_id":           validUUID,
			"origin_id":         validUUID,
			"sample_source_id":  validUUID,
			"microorganism_id":  validUUID,
			"sequencer_id":      validUUID,
			"laboratory_id":     validUUID,
			"health_service_id": validUUID,
		}

		c, w := testutils.SetupGinContext(http.MethodPost, "/api/admin/sample",
			testutils.ToJSON(minimalInput), nil, nil)
		handler.CreateSample(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Error - Invalid Gender", func(t *testing.T) {
		svc := &mocks.MockSampleService{}

		originalGender := validInput["gender"]
		validInput["gender"] = "monkey"
		defer func() {
			validInput["gender"] = originalGender
		}()

		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/sample",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)

		handler.CreateSample(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Invalid gender for sample.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Bad Request", func(t *testing.T) {
		svc := &mocks.MockSampleService{}
		handler := sample.NewAdminSampleHandler(svc)

		for _, test := range data.CreateSampleTests {
			t.Run(test.Name, func(t *testing.T) {
				c, w := testutils.SetupGinContext(
					http.MethodPost,
					"/api/admin/sample",
					test.Body,
					nil,
					nil,
				)

				handler.CreateSample(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, test.Expected, w.Body.String())
			})
		}
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			CreateFunc: func(ctx context.Context,
				input models.SampleCreateInput,
				language string) (*models.SampleResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/sample",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)

		handler.CreateSample(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
