package sample_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
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

	mockUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			FindAllFunc: func(ctx context.Context, input string,
				userID uuid.UUID, language string) (
				[]models.SampleResponse, error) {
				assert.Equal(t, mockUserID, userID)
				return []models.SampleResponse{mockResponse}, nil
			},
		}

		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/sample",
			"",
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
		handler.GetSamples(c)

		expected := testutils.ToJSON(
			map[string][]models.SampleResponse{
				"data": {mockResponse},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		svc := &mocks.MockSampleService{}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/sample",
			"",
			nil,
			nil,
		)
		handler.GetSamples(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			FindAllFunc: func(ctx context.Context, input string,
				userID uuid.UUID, language string) (
				[]models.SampleResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/sample",
			"",
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})
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
