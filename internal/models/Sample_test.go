package models_test

import (
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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
		"read1.fastq", "read2.fastq", "", mockCountry, mockUser, mockOrigin,
		mockSampleSource, mockMicro, mockSequencer, mockLab, mockHealthService,
	)
	language := "en"

	expected := models.SampleResponse{
		ID:             id,
		Name:           mockSample.Name,
		CollectionDate: date,
		RunNumber:      mockSample.RunNumber,
		RunDate:        mockSample.RunDate,
		City:           mockSample.City,
		OriginCode:     mockSample.OriginCode,
		Gender:         mockSample.Gender,
		DateOfBirth:    &date,
		Fastq1:         mockSample.Fastq1,
		Fastq2:         mockSample.Fastq2,
		Fasta:          mockSample.Fasta,
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
