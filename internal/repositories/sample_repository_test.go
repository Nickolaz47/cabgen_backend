package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createMockSample() models.Sample {
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

	return mockSample
}

func TestNewSampleRepo(t *testing.T) {
	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	assert.NotEmpty(t, sampleRepo)
}

func TestGetSamples(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := createMockSample()
	db.Create(&mockSample)

	t.Run("Success - All samples", func(t *testing.T) {
		result, err := sampleRepo.GetSamples(ctx, "")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockSample.ID, result[0].ID)
	})

	t.Run("Success - Filtered samples", func(t *testing.T) {
		result, err := sampleRepo.GetSamples(ctx, "neis")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockSample.ID, result[0].ID)
		assert.Equal(t, mockSample.Microorganism.Species,
			result[0].Microorganism.Species)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		samples, err := mockSampleRepo.GetSamples(
			context.Background(), "")

		assert.Empty(t, samples)
		assert.Error(t, err)
	})
}

func TestGetSampleByID(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := createMockSample()
	db.Create(&mockSample)

	t.Run("Success", func(t *testing.T) {
		resultSample, err := sampleRepo.GetSampleByID(ctx, mockSample.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockSample.ID, resultSample.ID)
		assert.Equal(t, mockSample.Name, resultSample.Name)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		resultSample, err := sampleRepo.GetSampleByID(ctx, uuid.New())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, resultSample)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		sample, err := mockSampleRepo.GetSampleByID(ctx, uuid.UUID{})

		assert.Empty(t, sample)
		assert.Error(t, err)
	})
}

func TestCreateSample(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := createMockSample()

	t.Run("Success", func(t *testing.T) {
		err := sampleRepo.CreateSample(ctx, &mockSample)
		assert.NoError(t, err)

		var result models.Sample
		err = db.Where("id = ?", mockSample.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, mockSample.ID, result.ID)
		assert.Equal(t, mockSample.Name, result.Name)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		err = mockSampleRepo.CreateSample(ctx, &models.Sample{})

		assert.Error(t, err)
	})
}

func TestUpdateSample(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := createMockSample()
	db.Create(&mockSample)

	t.Run("Success", func(t *testing.T) {
		sampleToUpdate := mockSample
		sampleToUpdate.Name = "Sample Updated Name"

		err := sampleRepo.UpdateSample(ctx, &sampleToUpdate)
		assert.NoError(t, err)

		var result models.Sample
		err = db.Where("id = ?", mockSample.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, "Sample Updated Name", result.Name)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		err = mockSampleRepo.UpdateSample(ctx, &models.Sample{})

		assert.Error(t, err)
	})
}

func TestDeleteSample(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := createMockSample()
	db.Create(&mockSample)

	t.Run("Success", func(t *testing.T) {
		err := sampleRepo.DeleteSample(ctx, &mockSample)
		assert.NoError(t, err)

		var result models.Sample
		err = db.Where("id = ?", mockSample.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		err = mockSampleRepo.DeleteSample(ctx, &models.Sample{})

		assert.Error(t, err)
	})
}
