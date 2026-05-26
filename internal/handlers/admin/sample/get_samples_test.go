package sample_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetSamples(t *testing.T) {
	testutils.SetupTestContext()

	mockSample := testmodels.CreateMockSample()
	mockResponse := mockSample.ToResponse("")

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			FindAllFunc: func(ctx context.Context, input string,
				userID uuid.UUID, language string) (
				[]models.SampleResponse, error) {
				return []models.SampleResponse{mockResponse}, nil
			},
		}

		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/sample",
			"",
			nil,
			nil,
		)

		handler.GetSamples(c)

		expected := testutils.ToJSON(
			map[string][]models.SampleResponse{
				"data": {mockResponse},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			FindAllFunc: func(ctx context.Context, input string,
				userID uuid.UUID, language string) (
				[]models.SampleResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/sample",
			"",
			nil,
			nil,
		)

		handler.GetSamples(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
