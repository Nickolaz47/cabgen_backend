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
		Languages: []models.SelectOption{
			{Label: "option.language.pt", Value: "pt"},
			{Label: "option.language.en", Value: "en"},
			{Label: "option.language.es", Value: "es"},
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

	labRepo := &mocks.MockLaboratoryRepository{
		GetActiveLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
			return []models.Laboratory{
				{ID: labID, Name: "LACEN/RJ"},
			}, nil
		},
	}
	seqRepo := &mocks.MockSequencerRepository{
		GetActiveSequencersFunc: func(ctx context.Context) ([]models.Sequencer, error) {
			return []models.Sequencer{
				{ID: seqID, Brand: "Illumina"},
			}, nil
		},
	}
	hsRepo := &mocks.MockHealthServiceRepository{
		GetActiveHealthServicesFunc: func(ctx context.Context) ([]models.HealthService, error) {
			return []models.HealthService{
				{ID: hsID, Name: "Hospital Central"},
			}, nil
		},
	}
	originRepo := &mocks.MockOriginRepository{
		GetActiveOriginsFunc: func(ctx context.Context) ([]models.Origin, error) {
			return []models.Origin{
				{ID: originID, Names: models.JSONMap{"en": "Human", "pt": "Humano"}},
			}, nil
		},
	}
	microRepo := &mocks.MockMicroorganismRepository{
		GetActiveMicroorganismsFunc: func(ctx context.Context) ([]models.Microorganism, error) {
			return []models.Microorganism{
				{ID: microID, Species: "Escherichia coli", Variety: models.JSONMap{"en": "", "pt": ""}},
			}, nil
		},
	}
	sourceRepo := &mocks.MockSampleSourceRepository{
		GetActiveSampleSourcesFunc: func(ctx context.Context) ([]models.SampleSource, error) {
			return []models.SampleSource{
				{ID: sourceID, Names: models.JSONMap{"en": "Aspirated", "pt": "Aspirado"}},
			}, nil
		},
	}

	expected := &models.FormSelectsResponse{
		Laboratories: []models.SelectOption{
			{Label: "LACEN/RJ", Value: labID.String()},
		},
		Sequencers: []models.SelectOption{
			{Label: "Illumina", Value: seqID.String()},
		},
		HealthServices: []models.SelectOption{
			{Label: "Hospital Central", Value: hsID.String()},
		},
		Origins: []models.SelectOption{
			{Label: "Humano", Value: originID.String()},
		},
		Microorganisms: []models.SelectOption{
			{Label: "Escherichia coli ", Value: microID.String()},
		},
		SampleSources: []models.SelectOption{
			{Label: "Aspirado", Value: sourceID.String()},
		},
	}

	t.Run("Success", func(t *testing.T) {
		svc := services.NewSelectOptionsService(
			labRepo, seqRepo, hsRepo, originRepo, microRepo, sourceRepo)

		result, err := svc.FindAllFormSelects(context.Background(), "pt")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
