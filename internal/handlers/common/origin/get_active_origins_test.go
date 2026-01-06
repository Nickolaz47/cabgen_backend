package origin

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockOriginService struct {
	FindAllFunc       func(ctx context.Context) ([]models.Origin, error)
	FindAllActiveFunc func(ctx context.Context, lang string) ([]models.OriginFormResponse, error)
	FindByIDFunc      func(ctx context.Context, ID uuid.UUID) (*models.Origin, error)
	FindByNameFunc    func(ctx context.Context, name, lang string) ([]models.Origin, error)
	CreateFunc        func(ctx context.Context, origin *models.Origin) error
	UpdateFunc        func(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error)
	DeleteFunc        func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockOriginService) FindAll(ctx context.Context) ([]models.Origin, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockOriginService) FindAllActive(ctx context.Context, lang string) ([]models.OriginFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx, lang)
	}
	return nil, nil
}

func (m *MockOriginService) FindByID(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (m *MockOriginService) FindByName(ctx context.Context, name, lang string) ([]models.Origin, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(ctx, name, lang)
	}
	return nil, nil
}

func (m *MockOriginService) Create(ctx context.Context, origin *models.Origin) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, origin)
	}
	return nil
}

func (m *MockOriginService) Update(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (m *MockOriginService) Delete(ctx context.Context, ID uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}
	return nil
}

func TestGetActiveOrigins(t *testing.T) {
	testutils.SetupTestContext()
	mockOrigin := models.Origin{
		ID:       uuid.New(),
		Names:    map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		originSvc := MockOriginService{
			FindAllActiveFunc: func(ctx context.Context, lang string) ([]models.OriginFormResponse, error) {
				return []models.OriginFormResponse{mockOrigin.ToFormResponse(lang)}, nil
			},
		}
		handler := NewOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/origin", "",
			nil, nil,
		)

		handler.GetActiveOrigins(c)

		expected := testutils.ToJSON(map[string][]models.OriginFormResponse{
			"data": {mockOrigin.ToFormResponse("en")},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		originSvc := MockOriginService{
			FindAllActiveFunc: func(ctx context.Context, lang string) ([]models.OriginFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := NewOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/origin", "",
			nil, nil,
		)

		handler.GetActiveOrigins(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
