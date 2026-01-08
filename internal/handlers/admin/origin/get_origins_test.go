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

	mockResponse := mockOrigin.ToAdminTableResponse("pt")

	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockOriginService{
			FindAllFunc: func(ctx context.Context, lang string) ([]models.OriginAdminTableResponse, error) {
				return []models.OriginAdminTableResponse{mockResponse}, nil
			},
		}

		handler := origin.NewAdminOriginHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin",
			"",
			nil,
			nil,
		)

		handler.GetAllOrigins(c)

		expected := testutils.ToJSON(
			map[string][]models.OriginAdminTableResponse{
				"data": {mockResponse},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := testmodels.MockOriginService{
			FindAllFunc: func(ctx context.Context, lang string) ([]models.OriginAdminTableResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := origin.NewAdminOriginHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin",
			"",
			nil,
			nil,
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
