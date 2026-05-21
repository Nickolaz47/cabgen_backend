package models_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToGender(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected models.Gender
	}{
		{
			name:     "String to Male",
			input:    "masculino",
			expected: models.Male,
		},
		{
			name:     "String to Female",
			input:    "femenino",
			expected: models.Female,
		},
		{
			name:     "String to Unspecified",
			input:    "any string",
			expected: models.Unspecified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.ToGender(tt.input)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToTranslatedString(t *testing.T) {
	tests := []struct {
		name     string
		language string
		gender   models.Gender
		expected string
	}{
		{
			name:     "Male to portuguese",
			language: "pt",
			gender:   models.Male,
			expected: "Masculino",
		},
		{
			name:     "Female to spanish",
			language: "es",
			gender:   models.Female,
			expected: "Femenino",
		},
		{
			name:     "Unspecified to english",
			language: "en",
			gender:   models.Unspecified,
			expected: "Unspecified",
		},
		{
			name:     "Invalid language",
			language: "an",
			gender:   models.Male,
			expected: "Male",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.gender.ToTranslatedString(tt.language)
			assert.Equal(t, &tt.expected, result)
		})
	}

	t.Run("Gender is nil", func(t *testing.T) {
		var gender *models.Gender
		result := gender.ToTranslatedString("en")

		assert.Nil(t, result)
	})
}

func TestSampleToResponse(t *testing.T) {
	mockUser := testmodels.NewLoginUser()
	mockCountry := mockUser.Country
	mockOrigin := testmodels.NewOrigin(uuid.New().String(),
		map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		true,
	)
	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado",
			"en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório",
			"en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)
	mockMicro := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B"}, true,
	)
	mockSequencer := models.Sequencer{
		ID:       uuid.New(),
		Brand:    "Illumina",
		Model:    "MySeq",
		IsActive: true,
	}
	mockLab := models.Laboratory{
		ID:           uuid.New(),
		Name:         "Laboratorio Central do Rio de Janeiro",
		Abbreviation: "LACEN/RJ",
		IsActive:     true,
	}
	mockHealthService := models.HealthService{
		ID:           uuid.New(),
		Name:         "Laboratorio Central do Rio de Janeiro",
		Type:         "Public",
		CountryID:    mockCountry.ID,
		Country:      mockCountry,
		City:         nil,
		Contactant:   nil,
		ContactEmail: nil,
		ContactPhone: nil,
		IsActive:     true,
	}

	id := uuid.New()
	date := time.Date(2024, time.May, 11, 0, 0, 0, 0, time.UTC)

	mockSample := testmodels.NewSample(
		id.String(), "sample 1", date, "R1", date, "", "A01", models.Male, date,
		"sequences/read1.fastq", "sequences/read2.fastq",
		"sequences/read.fasta", mockCountry,
		mockUser, mockOrigin, mockSampleSource, mockMicro, mockSequencer,
		mockLab, mockHealthService,
	)
	language := "en"

	fastq1 := filepath.Base(*mockSample.Fastq1)
	fastq2 := filepath.Base(*mockSample.Fastq2)
	fasta := filepath.Base(*mockSample.Fasta)

	expected := models.SampleResponse{
		ID:             id,
		Name:           mockSample.Name,
		CollectionDate: date,
		RunNumber:      mockSample.RunNumber,
		RunDate:        mockSample.RunDate,
		City:           mockSample.City,
		OriginCode:     mockSample.OriginCode,
		Gender:         mockSample.Gender.ToTranslatedString(language),
		DateOfBirth:    &date,
		Fastq1:         &fastq1,
		Fastq2:         &fastq2,
		Fasta:          &fasta,
		CountryCode:    mockCountry.Code,
		User:           mockUser.Username,
		Origin:         mockOrigin.Names[language],
		SampleSource:   mockSampleSource.Names[language],
		Microorganism:  mockMicro.Species + " " + mockMicro.Variety[language],
		Sequencer:      mockSequencer.Brand + " - " + mockSequencer.Model,
		Laboratory:     mockLab.Name,
		HealthService:  mockHealthService.Name,
	}
	result := mockSample.ToResponse("")

	assert.Equal(t, expected, result)
}
