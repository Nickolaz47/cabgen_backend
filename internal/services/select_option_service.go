package services

import (
	"context"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type SelectOptionsService interface {
	FindAll(ctx context.Context) (*models.EnumSelectsResponse, error)
}

type selectOptionsService struct{}

func NewSelectOptionsService() SelectOptionsService {
	return &selectOptionsService{}
}

func (s *selectOptionsService) FindAll(ctx context.Context) (
	*models.EnumSelectsResponse, error) {
	resp := &models.EnumSelectsResponse{}

	// User Roles
	for _, role := range models.UserRoles {
		resp.Roles = append(resp.Roles, models.SelectOption{
			Label: "option.role." + strings.ToLower(string(role)),
			Value: string(role),
		})
	}

	// Taxons
	for _, taxon := range models.Taxons {
		resp.Taxons = append(resp.Taxons, models.SelectOption{
			Label: "option.taxon." + strings.ToLower(string(taxon)),
			Value: string(taxon),
		})
	}

	// Genders
	for _, gender := range models.Genders {
		resp.Genders = append(resp.Genders, models.SelectOption{
			Label: "option.gender." + strings.ToLower(string(gender)),
			Value: string(gender),
		})
	}

	// Health Service Types
	for _, hsType := range models.HealthServiceTypes {
		resp.HealthServiceTypes = append(resp.HealthServiceTypes,
			models.SelectOption{
				Label: "option.health_service_type." + strings.ToLower(
					string(hsType)),
				Value: string(hsType),
			})
	}

	// Analisis Types
	for _, aType := range models.AnalysisTypes {
		resp.AnalysisTypes = append(resp.AnalysisTypes, models.SelectOption{
			Label: "option.analysis_type." + strings.ToLower(string(aType)),
			Value: string(aType),
		})
	}

	return resp, nil
}
