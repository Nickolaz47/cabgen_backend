package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/xuri/excelize/v2"
)

type MockTemplateService struct {
	BuildColumnsFunc func(ctx context.Context, language string) (
		[]utils.Column, map[string]struct{}, error)
	CreateTemplateTableFunc func(ctx context.Context,
		language string) (*excelize.File, error)
	ValidateTemplateTableFunc func(ctx context.Context, language string,
		file *excelize.File) ([]models.SampleCreateInput, error)
}

func (m *MockTemplateService) BuildColumns(ctx context.Context,
	language string) ([]utils.Column, map[string]struct{}, error) {
	if m.BuildColumnsFunc != nil {
		return m.BuildColumnsFunc(ctx, language)
	}

	return nil, nil, nil
}

func (m *MockTemplateService) CreateTemplateTable(ctx context.Context,
	language string) (*excelize.File, error) {
	if m.CreateTemplateTableFunc != nil {
		return m.CreateTemplateTableFunc(ctx, language)
	}

	return nil, nil
}

func (m *MockTemplateService) ValidateTemplateTable(ctx context.Context,
	language string, file *excelize.File) ([]models.SampleCreateInput,
	error) {
	if m.ValidateTemplateTableFunc != nil {
		return m.ValidateTemplateTableFunc(ctx, language, file)
	}

	return nil, nil
}
