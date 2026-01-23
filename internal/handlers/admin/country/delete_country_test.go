package country_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/country"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCountry(t *testing.T) {
	testutils.SetupTestContext()

	mockCountry := testmodels.NewCountry("", nil)
	mockCountry.ID = 1

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			DeleteFunc: func(ctx context.Context, code string) error {
				return nil
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/country",
			"", nil,
			gin.Params{{Key: "code", Value: mockCountry.Code}},
		)
		handler.DeleteCountry(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "Country deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			DeleteFunc: func(ctx context.Context, code string) error {
				return services.ErrNotFound
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/country",
			"", nil, gin.Params{{Key: "code", Value: "ARG"}},
		)
		handler.DeleteCountry(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "No country was found with this code.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			DeleteFunc: func(ctx context.Context, code string) error {
				return services.ErrInternal
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/country",
			"", nil, gin.Params{{Key: "code", Value: mockCountry.Code}},
		)
		handler.DeleteCountry(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
