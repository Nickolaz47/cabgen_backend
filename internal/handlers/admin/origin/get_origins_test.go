package origin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllOrigins(t *testing.T) {
	testutils.SetupTestContext()
	mockOrigin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		originSvc := testmodels.MockOriginService{
			FindAllFunc: func(ctx context.Context) ([]models.Origin, error) {
				return []models.Origin{mockOrigin}, nil
			},
		}

		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin", "",
			nil, nil,
		)
		handler.GetAllOrigins(c)

		expected := testutils.ToJSON(
			map[string][]models.Origin{
				"data": {mockOrigin},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		originSvc := testmodels.MockOriginService{
			FindAllFunc: func(ctx context.Context) ([]models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/origin", "",
			nil, nil,
		)
		handler.GetAllOrigins(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
