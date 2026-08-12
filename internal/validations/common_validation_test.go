package validations_test

import (
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestApplySampleUpdate(t *testing.T) {
	mock := testmodels.CreateMockSample()

	collectionDate := models.Date{
		Time: time.Date(2025, time.January, 10, 0, 0, 0, 0, time.UTC)}
	runNumber := "R2"
	runDate := models.Date{
		Time: time.Date(2025, time.February, 5, 0, 0, 0, 0, time.UTC)}
	city := "Rio de Janeiro"
	originCode := "B02"
	gender := models.Female
	dateOfBirth := models.Date{
		Time: time.Date(1990, time.March, 15, 0, 0, 0, 0, time.UTC)}

	input := models.SampleUpdateDTO{
		OriginCode:     &originCode,
		CollectionDate: &collectionDate,
		RunNumber:      &runNumber,
		RunDate:        &runDate,
		City:           &city,
		Gender:         &gender,
		DateOfBirth:    &dateOfBirth,
	}

	expected := models.Sample{
		ID:              mock.ID,
		OriginCode:      originCode,
		CollectionDate:  *models.ToTimePtr(&collectionDate),
		RunNumber:       runNumber,
		RunDate:         *models.ToTimePtr(&runDate),
		City:            &city,
		Gender:          &gender,
		DateOfBirth:     models.ToTimePtr(&dateOfBirth),
		CountryID:       mock.CountryID,
		Country:         mock.Country,
		UserID:          mock.UserID,
		User:            mock.User,
		OriginID:        mock.OriginID,
		Origin:          mock.Origin,
		SampleSourceID:  mock.SampleSourceID,
		SampleSource:    mock.SampleSource,
		MicroorganismID: mock.MicroorganismID,
		Microorganism:   mock.Microorganism,
		SequencerID:     mock.SequencerID,
		Sequencer:       mock.Sequencer,
		LaboratoryID:    mock.LaboratoryID,
		Laboratory:      mock.Laboratory,
		HealthServiceID: mock.HealthServiceID,
		HealthService:   mock.HealthService,
		Fastq1:          mock.Fastq1,
		Fastq2:          mock.Fastq2,
		Fasta:           mock.Fasta,
	}

	validations.ApplySampleUpdate(&mock, &input)

	assert.Equal(t, expected, mock)
}

func TestApplySampleFilesUpdate(t *testing.T) {
	mock := testmodels.CreateMockSample()

	fastq1 := "new_read1.fastq"
	fastq2 := "new_read2.fastq"
	fasta := "new_assembly.fasta"

	input := models.SampleAttachmentInput{
		Fastq1: &fastq1,
		Fastq2: &fastq2,
		Fasta:  &fasta,
	}

	expected := models.Sample{
		ID:              mock.ID,
		OriginCode:      mock.OriginCode,
		CollectionDate:  mock.CollectionDate,
		RunNumber:       mock.RunNumber,
		RunDate:         mock.RunDate,
		City:            mock.City,
		Gender:          mock.Gender,
		DateOfBirth:     mock.DateOfBirth,
		CountryID:       mock.CountryID,
		Country:         mock.Country,
		UserID:          mock.UserID,
		User:            mock.User,
		OriginID:        mock.OriginID,
		Origin:          mock.Origin,
		SampleSourceID:  mock.SampleSourceID,
		SampleSource:    mock.SampleSource,
		MicroorganismID: mock.MicroorganismID,
		Microorganism:   mock.Microorganism,
		SequencerID:     mock.SequencerID,
		Sequencer:       mock.Sequencer,
		LaboratoryID:    mock.LaboratoryID,
		Laboratory:      mock.Laboratory,
		HealthServiceID: mock.HealthServiceID,
		HealthService:   mock.HealthService,
		Fastq1:          &fastq1,
		Fastq2:          &fastq2,
		Fasta:           &fasta,
	}

	validations.ApplySampleFilesUpdate(&mock, &input)

	assert.Equal(t, expected, mock)
}
