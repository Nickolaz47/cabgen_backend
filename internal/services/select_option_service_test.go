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
	tests := []struct {
		language string
		expected models.EnumSelectsResponse
	}{
		{
			language: "pt",
			expected: models.EnumSelectsResponse{
				Roles: []models.SelectOption{
					{Label: "Administrador", Value: "Admin"},
					{Label: "Colaborador", Value: "Collaborator"},
				},
				Taxons: []models.SelectOption{
					{Label: "Bactéria", Value: "Bacteria"},
					{Label: "Fungos", Value: "Fungi"},
					{Label: "Protozoa", Value: "Protozoa"},
					{Label: "Vírus", Value: "Virus"},
				},
				Genders: []models.SelectOption{
					{Label: "Feminino", Value: "Female"},
					{Label: "Masculino", Value: "Male"},
					{Label: "Não especificado", Value: "Unspecified"},
				},
				HealthServiceTypes: []models.SelectOption{
					{Label: "Público", Value: "Public"},
					{Label: "Privado", Value: "Private"},
				},
				AnalysisTypes: []models.SelectOption{
					{Label: "Qualidade", Value: "FASTQC"},
					{Label: "Genômica", Value: "GENOME"},
					{Label: "Completa", Value: "COMPLETE"},
				},
				Languages: []models.SelectOption{
					{Label: "Português", Value: "pt"},
					{Label: "Inglês", Value: "en"},
					{Label: "Espanhol", Value: "es"},
				},
			},
		},
		{
			language: "en",
			expected: models.EnumSelectsResponse{
				Roles: []models.SelectOption{
					{Label: "Admin", Value: "Admin"},
					{Label: "Collaborator", Value: "Collaborator"},
				},
				Taxons: []models.SelectOption{
					{Label: "Bacteria", Value: "Bacteria"},
					{Label: "Fungi", Value: "Fungi"},
					{Label: "Protozoa", Value: "Protozoa"},
					{Label: "Virus", Value: "Virus"},
				},
				Genders: []models.SelectOption{
					{Label: "Female", Value: "Female"},
					{Label: "Male", Value: "Male"},
					{Label: "Unspecified", Value: "Unspecified"},
				},
				HealthServiceTypes: []models.SelectOption{
					{Label: "Public", Value: "Public"},
					{Label: "Private", Value: "Private"},
				},
				AnalysisTypes: []models.SelectOption{
					{Label: "Quality", Value: "FASTQC"},
					{Label: "Genomic", Value: "GENOME"},
					{Label: "Complete", Value: "COMPLETE"},
				},
				Languages: []models.SelectOption{
					{Label: "Portuguese", Value: "pt"},
					{Label: "English", Value: "en"},
					{Label: "Spanish", Value: "es"},
				},
			},
		},
		{
			language: "es",
			expected: models.EnumSelectsResponse{
				Roles: []models.SelectOption{
					{Label: "Administrador", Value: "Admin"},
					{Label: "Colaborador", Value: "Collaborator"},
				},
				Taxons: []models.SelectOption{
					{Label: "Bacteria", Value: "Bacteria"},
					{Label: "Hongos", Value: "Fungi"},
					{Label: "Protozoa", Value: "Protozoa"},
					{Label: "Virus", Value: "Virus"},
				},
				Genders: []models.SelectOption{
					{Label: "Femenino", Value: "Female"},
					{Label: "Masculino", Value: "Male"},
					{Label: "No especificado", Value: "Unspecified"},
				},
				HealthServiceTypes: []models.SelectOption{
					{Label: "Público", Value: "Public"},
					{Label: "Privado", Value: "Private"},
				},
				AnalysisTypes: []models.SelectOption{
					{Label: "Calidad", Value: "FASTQC"},
					{Label: "Genómica", Value: "GENOME"},
					{Label: "Completo", Value: "COMPLETE"},
				},
				Languages: []models.SelectOption{
					{Label: "Portugués", Value: "pt"},
					{Label: "Inglés", Value: "en"},
					{Label: "Español", Value: "es"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			svc := services.NewSelectOptionsService(nil, nil, nil, nil,
				nil, nil)

			result, err := svc.FindAllEnumSelects(context.Background(),
				tt.language)

			assert.NoError(t, err)
			assert.Equal(t, &tt.expected, result)
		})
	}
}

func TestSelectOptionFindAllFormSelects(t *testing.T) {
	labID := uuid.New()
	otherLabID := uuid.New()
	seqID := uuid.New()
	hsID := uuid.New()
	originID := uuid.New()
	microID := uuid.New()
	sourceID := uuid.New()

	labRepo := &mocks.MockLaboratoryRepository{
		GetActiveLaboratoriesFunc: func(ctx context.Context) ([]models.Laboratory, error) {
			return []models.Laboratory{
				{ID: labID, Name: "LACEN/RJ"},
				{ID: otherLabID, Name: "option.laboratory.other"},
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
			{Label: "Outro", Value: otherLabID.String()},
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
		Genders: []models.SelectOption{
			{Label: "Feminino", Value: "Female"},
			{Label: "Masculino", Value: "Male"},
			{Label: "Não especificado", Value: "Unspecified"},
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
