package laboratory_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockLaboratoryService struct {
	FindAllFunc                  func(ctx context.Context) ([]models.Laboratory, error)
	FindAllActiveFunc            func(ctx context.Context) ([]models.LaboratoryFormResponse, error)
	FindByIDFunc                 func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error)
	FindByNameOrAbbreviationFunc func(ctx context.Context, input string) ([]models.Laboratory, error)
	CreateFunc                   func(ctx context.Context, lab *models.Laboratory) error
	UpdateFunc                   func(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error)
	DeleteFunc                   func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockLaboratoryService) FindAll(ctx context.Context) ([]models.Laboratory, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindAllActive(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error) {
	if m.FindByNameOrAbbreviationFunc != nil {
		return m.FindByNameOrAbbreviationFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockLaboratoryService) Create(ctx context.Context, lab *models.Laboratory) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, lab)
	}
	return nil
}

func (m *MockLaboratoryService) Update(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (m *MockLaboratoryService) Delete(ctx context.Context, ID uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}
	return nil
}

func TestGetActiveLaboratories(t *testing.T) {
	testutils.SetupTestContext()

	mockLab := models.LaboratoryFormResponse{ID: uuid.New()}

	t.Run("Success", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			FindAllActiveFunc: func(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
				return []models.LaboratoryFormResponse{mockLab}, nil
			},
		}

		handler := laboratory.NewLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/laboratory", "",
			nil, nil,
		)
		handler.GetActiveLaboratories(c)

		expected := testutils.ToJSON(
			map[string][]models.LaboratoryFormResponse{
				"data": {mockLab},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			FindAllActiveFunc: func(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := laboratory.NewLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/laboratory", "",
			nil, nil,
		)
		handler.GetActiveLaboratories(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
