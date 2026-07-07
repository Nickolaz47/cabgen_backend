package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestSelectOptionFindAll(t *testing.T) {
	svc := services.NewSelectOptionsService()

	expected := models.EnumSelectsResponse{
		Roles: []models.SelectOption{
			{Label: "option.role.admin", Value: "Admin"},
			{Label: "option.role.collaborator", Value: "Collaborator"},
		},
		Taxons: []models.SelectOption{
			{Label: "option.taxon.bacteria", Value: "Bacteria"},
			{Label: "option.taxon.fungi", Value: "Fungi"},
			{Label: "option.taxon.protozoa", Value: "Protozoa"},
			{Label: "option.taxon.virus", Value: "Virus"},
		},
		Genders: []models.SelectOption{
			{Label: "option.gender.female", Value: "Female"},
			{Label: "option.gender.male", Value: "Male"},
			{Label: "option.gender.unspecified", Value: "Unspecified"},
		},
		HealthServiceTypes: []models.SelectOption{
			{Label: "option.health_service_type.public", Value: "Public"},
			{Label: "option.health_service_type.private", Value: "Private"},
		},
		AnalysisTypes: []models.SelectOption{
			{Label: "option.analysis_type.fastqc", Value: "FASTQC"},
			{Label: "option.analysis_type.genome", Value: "GENOME"},
			{Label: "option.analysis_type.complete", Value: "COMPLETE"},
		},
	}

	result, err := svc.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, &expected, result)
}
