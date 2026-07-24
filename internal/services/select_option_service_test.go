package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSelectOptionFindAllEnumSelects(t *testing.T) {
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

	t.Run("Success", func(t *testing.T) {
		svc := services.NewSelectOptionsService(nil, nil, nil, nil, nil, nil)

		result, err := svc.FindAllEnumSelects(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, &expected, result)
	})
}

func TestSelectOptionFindAllFormSelects(t *testing.T) {
	labID := uuid.New()
	seqID := uuid.New()
	hsID := uuid.New()
	originID := uuid.New()
	microID := uuid.New()
	sourceID := uuid.New()

	labSvc := &mocks.MockLaboratoryService{
		FindAllActiveFunc: func(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
			return []models.LaboratoryFormResponse{
				{ID: labID, Name: "LACEN/RJ"},
			}, nil
		},
	}
	seqSvc := &mocks.MockSequencerService{
		FindAllActiveFunc: func(ctx context.Context) ([]models.SequencerFormResponse, error) {
			return []models.SequencerFormResponse{
				{ID: seqID, Model: "MiSeq"},
			}, nil
		},
	}
	hsSvc := &mocks.MockHealthServiceService{
		FindAllActiveFunc: func(ctx context.Context) ([]models.HealthServiceFormResponse, error) {
			return []models.HealthServiceFormResponse{
				{ID: hsID, Name: "Hospital Central"},
			}, nil
		},
	}
	originSvc := &mocks.MockOriginService{
		FindAllActiveFunc: func(ctx context.Context, language string) ([]models.OriginFormResponse, error) {
			return []models.OriginFormResponse{
				{ID: originID, Name: "Humano"},
			}, nil
		},
	}
	microSvc := &mocks.MockMicroorganismService{
		FindAllActiveFunc: func(ctx context.Context, language string) ([]models.MicroorganismFormResponse, error) {
			return []models.MicroorganismFormResponse{
				{ID: microID, Species: "Escherichia coli"},
			}, nil
		},
	}
	sourceSvc := &mocks.MockSampleSourceService{
		FindAllActiveFunc: func(ctx context.Context, language string) ([]models.SampleSourceFormResponse, error) {
			return []models.SampleSourceFormResponse{
				{ID: sourceID, Name: "Aspirado"},
			}, nil
		},
	}

	expected := &models.FormSelectsResponse{
		Laboratories: []models.SelectOption{
			{Label: "LACEN/RJ", Value: labID.String()},
		},
		Sequencers: []models.SelectOption{
			{Label: "MiSeq", Value: seqID.String()},
		},
		HealthServices: []models.SelectOption{
			{Label: "Hospital Central", Value: hsID.String()},
		},
		Origins: []models.SelectOption{
			{Label: "Humano", Value: originID.String()},
		},
		Microorganisms: []models.SelectOption{
			{Label: "Escherichia coli", Value: microID.String()},
		},
		SampleSources: []models.SelectOption{
			{Label: "Aspirado", Value: sourceID.String()},
		},
	}

	t.Run("Success", func(t *testing.T) {
		svc := services.NewSelectOptionsService(
			labSvc, seqSvc, hsSvc, originSvc, microSvc, sourceSvc)

		result, err := svc.FindAllFormSelects(context.Background(), "pt")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
