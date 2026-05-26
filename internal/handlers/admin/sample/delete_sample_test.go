package sample_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSample(t *testing.T) {
	testutils.SetupTestContext()

	mockSample := testmodels.CreateMockSample()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			DeleteFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID) error {
				return nil
			},
		}
		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/sample",
			"",
			nil,
			gin.Params{{Key: "sampleId", Value: mockSample.ID.String()}},
		)

		handler.DeleteSample(c)

		expected := testutils.ToJSON(
			map[string]any{
				"message": "Sample deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockSampleService{}
		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/sample",
			"",
			nil,
			nil,
		)

		handler.DeleteSample(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not found", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			DeleteFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID) error {
				return services.ErrNotFound
			},
		}

		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/sample",
			"",
			nil,
			gin.Params{{Key: "sampleId", Value: mockSample.ID.String()}},
		)

		handler.DeleteSample(c)

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
			DeleteFunc: func(ctx context.Context, sampleID,
				userID uuid.UUID) error {
				return services.ErrInternal
			},
		}

		handler := sample.NewAdminSampleHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/sample",
			"",
			nil,
			gin.Params{{Key: "sampleId", Value: mockSample.ID.String()}},
		)

		handler.DeleteSample(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
