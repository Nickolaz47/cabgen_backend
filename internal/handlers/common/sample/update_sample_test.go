package sample_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSample(t *testing.T) {
	testutils.SetupTestContext()

	mockSample := testmodels.CreateMockSample()
	mockResponse := mockSample.ToResponse("")

	validUpdateInput := map[string]any{
		"name":       "Updated-Sample-Name",
		"run_number": "RUN-UPDATED-01",
		"city":       "Niterói",
		"gender":     "Female",
	}

	mockUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			UpdateFunc: func(ctx context.Context, sampleID, userID uuid.UUID,
				input models.SampleUpdateDTO, language string) (
				*models.SampleResponse, error) {
				assert.Equal(t, mockUserID, userID)
				assert.NotNil(t, input.UserID)
				assert.Equal(t, mockUserID, *input.UserID)
				return &mockResponse, nil
			},
		}

		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/sample",
			testutils.ToJSON(validUpdateInput),
			nil,
			gin.Params{{Key: "sampleId", Value: mockSample.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

		handler.UpdateSample(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockResponse,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockSampleService{}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/sample",
			"",
			nil,
			nil,
		)

		handler.UpdateSample(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		svc := &mocks.MockSampleService{}
		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/sample",
			testutils.ToJSON(validUpdateInput),
			nil,
			gin.Params{{Key: "sampleId", Value: mockSample.ID.String()}},
		)

		handler.UpdateSample(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Bad Request", func(t *testing.T) {
		svc := &mocks.MockSampleService{}
		handler := sample.NewSampleHandler(svc)

		for _, test := range data.UpdateSampleTests {
			t.Run(test.Name, func(t *testing.T) {
				c, w := testutils.SetupGinContext(
					http.MethodPut,
					"/api/sample",
					test.Body,
					nil,
					gin.Params{{Key: "sampleId", Value: mockSample.ID.String()}},
				)

				c.Set("user", &models.UserToken{ID: mockUserID})

				handler.UpdateSample(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, test.Expected, w.Body.String())
			})
		}
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			UpdateFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleUpdateDTO,
				language string) (*models.SampleResponse, error) {
				return nil, services.ErrNotFound
			},
		}

		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/sample",
			testutils.ToJSON(validUpdateInput),
			nil,
			gin.Params{{Key: "sampleId", Value: uuid.NewString()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

		handler.UpdateSample(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Sample not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			UpdateFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID, input models.SampleUpdateDTO,
				language string) (*models.SampleResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := sample.NewSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/sample",
			testutils.ToJSON(validUpdateInput),
			nil,
			gin.Params{{Key: "sampleId", Value: mockSample.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

		handler.UpdateSample(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
