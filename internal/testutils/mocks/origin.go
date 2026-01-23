package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockOriginRepository struct {
	GetOriginsFunc         func(ctx context.Context) ([]models.Origin, error)
	GetActiveOriginsFunc   func(ctx context.Context) ([]models.Origin, error)
	GetOriginByIDFunc      func(ctx context.Context, ID uuid.UUID) (*models.Origin, error)
	GetOriginsByNameFunc   func(ctx context.Context, name, lang string) ([]models.Origin, error)
	GetOriginDuplicateFunc func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error)
	CreateOriginFunc       func(ctx context.Context, origin *models.Origin) error
	UpdateOriginFunc       func(ctx context.Context, origin *models.Origin) error
	DeleteOriginFunc       func(ctx context.Context, origin *models.Origin) error
}

func (r *MockOriginRepository) GetOrigins(ctx context.Context) ([]models.Origin, error) {
	if r.GetOriginsFunc != nil {
		return r.GetOriginsFunc(ctx)
	}
	return nil, nil
}

func (r *MockOriginRepository) GetActiveOrigins(ctx context.Context) ([]models.Origin, error) {
	if r.GetActiveOriginsFunc != nil {
		return r.GetActiveOriginsFunc(ctx)
	}
	return nil, nil
}

func (r *MockOriginRepository) GetOriginByID(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
	if r.GetOriginByIDFunc != nil {
		return r.GetOriginByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (r *MockOriginRepository) GetOriginsByName(ctx context.Context, name, lang string) ([]models.Origin, error) {
	if r.GetOriginsByNameFunc != nil {
		return r.GetOriginsByNameFunc(ctx, name, lang)
	}
	return nil, nil
}

func (r *MockOriginRepository) GetOriginDuplicate(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
	if r.GetOriginDuplicateFunc != nil {
		return r.GetOriginDuplicateFunc(ctx, names, ID)
	}
	return nil, nil
}

func (r *MockOriginRepository) CreateOrigin(ctx context.Context, origin *models.Origin) error {
	if r.CreateOriginFunc != nil {
		return r.CreateOriginFunc(ctx, origin)
	}
	return nil
}

func (r *MockOriginRepository) UpdateOrigin(ctx context.Context, origin *models.Origin) error {
	if r.UpdateOriginFunc != nil {
		return r.UpdateOriginFunc(ctx, origin)
	}
	return nil
}

func (r *MockOriginRepository) DeleteOrigin(ctx context.Context, origin *models.Origin) error {
	if r.DeleteOriginFunc != nil {
		return r.DeleteOriginFunc(ctx, origin)
	}
	return nil
}

type MockOriginService struct {
	FindAllFunc       func(ctx context.Context, lang string) ([]models.OriginAdminTableResponse, error)
	FindAllActiveFunc func(ctx context.Context, lang string) ([]models.OriginFormResponse, error)
	FindByIDFunc      func(ctx context.Context, ID uuid.UUID) (*models.OriginAdminDetailResponse, error)
	FindByNameFunc    func(ctx context.Context, name, lang string) ([]models.OriginAdminTableResponse, error)
	CreateFunc        func(ctx context.Context, input models.OriginCreateInput) (*models.OriginAdminDetailResponse, error)
	UpdateFunc        func(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.OriginAdminDetailResponse, error)
	DeleteFunc        func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockOriginService) FindAll(ctx context.Context, lang string) ([]models.OriginAdminTableResponse, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, lang)
	}
	return nil, nil
}

func (m *MockOriginService) FindAllActive(ctx context.Context, lang string) ([]models.OriginFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx, lang)
	}
	return nil, nil
}

func (m *MockOriginService) FindByID(ctx context.Context, ID uuid.UUID) (*models.OriginAdminDetailResponse, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (m *MockOriginService) FindByName(ctx context.Context, name, lang string) ([]models.OriginAdminTableResponse, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(ctx, name, lang)
	}
	return nil, nil
}

func (m *MockOriginService) Create(ctx context.Context, input models.OriginCreateInput) (*models.OriginAdminDetailResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockOriginService) Update(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.OriginAdminDetailResponse, error) {
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
