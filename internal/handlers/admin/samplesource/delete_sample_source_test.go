package samplesource_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSampleSource(t *testing.T) {
	testutils.SetupTestContext()
	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return nil
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sample-source", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)
		handler.DeleteSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "Sample source deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Invalid ID", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sample-source", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: "12"}},
		)
		handler.DeleteSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Sample source not found", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrNotFound
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sample-source", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)
		handler.DeleteSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Sample source not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("DB error", func(t *testing.T) {
		svc := testmodels.MockSampleSourceService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrInternal
			},
		}
		handler := samplesource.NewAdminSampleSourceHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sample-source", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)
		handler.DeleteSampleSource(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
