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
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
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
				userID uuid.UUID, offset int, language string) (
				[]models.SampleResponse, int, error) {
				assert.Equal(t, mockUserID, userID)
				return []models.SampleResponse{mockResponse}, 1, nil
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
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})
		handler.GetSamples(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data":        []models.SampleResponse{mockResponse},
				"total_pages": 1,
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
				userID uuid.UUID, offset int, language string) (
				[]models.SampleResponse, int, error) {
				return nil, 0, services.ErrInternal
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
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})
		handler.GetSamples(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Page", func(t *testing.T) {
		svc := &mocks.MockSampleService{}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/sample?page=0",
			"",
			nil,
			nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})
		handler.GetSamples(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Invalid query parameters.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - With page", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			FindAllFunc: func(ctx context.Context, input string,
				userID uuid.UUID, offset int, language string) (
				[]models.SampleResponse, int, error) {
				assert.Equal(t, 100, offset)
				return []models.SampleResponse{mockResponse}, 5, nil
			},
		}

		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/sample?page=2",
			"",
			nil,
			nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})
		handler.GetSamples(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data":        []models.SampleResponse{mockResponse},
				"total_pages": 5,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

}
