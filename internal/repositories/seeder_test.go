package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCountryBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	country1 := testmodels.NewCountry("BRA", map[string]string{
		"pt": "Brasil", "en": "Brazil", "es": "Brazil",
	})
	country2 := testmodels.NewCountry("SPN", map[string]string{
		"pt": "Espanha", "en": "Spain", "es": "España",
	})
	countries := []models.Country{country1, country2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewCountrySeedRepository(db)

		err := repo.BulkInsert(context.Background(), countries)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewCountrySeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), countries)

		assert.Error(t, err)
	})
}

func TestCountryCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewCountrySeedRepository(db)

	country1 := testmodels.NewCountry("BRA", map[string]string{
		"pt": "Brasil", "en": "Brazil", "es": "Brazil",
	})
	country2 := testmodels.NewCountry("SPN", map[string]string{
		"pt": "Espanha", "en": "Spain", "es": "España",
	})
	db.Create(&country1)
	db.Create(&country2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewCountrySeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}

func TestMicroorganismBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	micro1 := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Escherichia coli",
		map[string]string{
			"pt": "Enteropatogênica A",
			"en": "Enteropathogenic A",
			"es": "Enteropatogénica A",
		}, true,
	)
	micro2 := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Plesiomonas shigelloides",
		nil, true,
	)
	micros := []models.Microorganism{micro1, micro2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewMicroorganismSeedRepository(db)

		err := repo.BulkInsert(context.Background(), micros)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewMicroorganismSeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), micros)

		assert.Error(t, err)
	})
}

func TestMicroorganismCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismSeedRepository(db)

	micro1 := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Escherichia coli",
		map[string]string{
			"pt": "Enteropatogênica A",
			"en": "Enteropathogenic A",
			"es": "Enteropatogénica A",
		}, true,
	)
	micro2 := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Plesiomonas shigelloides",
		nil, true,
	)
	db.Create(&micro1)
	db.Create(&micro2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismSeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}

func TestOriginBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	origin1 := testmodels.NewOrigin(uuid.NewString(), map[string]string{
		"pt": "Humano", "en": "Human", "es": "Humano",
	}, true)
	origin2 := testmodels.NewOrigin(uuid.NewString(), map[string]string{
		"pt": "Animal", "en": "Animal", "es": "Animal",
	}, true)
	origins := []models.Origin{origin1, origin2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewOriginSeedRepository(db)

		err := repo.BulkInsert(context.Background(), origins)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewOriginSeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), origins)

		assert.Error(t, err)
	})
}

func TestOriginCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewOriginSeedRepository(db)

	origin1 := testmodels.NewOrigin(uuid.NewString(), map[string]string{
		"pt": "Humano", "en": "Human", "es": "Humano",
	}, true)
	origin2 := testmodels.NewOrigin(uuid.NewString(), map[string]string{
		"pt": "Animal", "en": "Animal", "es": "Animal",
	}, true)
	db.Create(&origin1)
	db.Create(&origin2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewOriginSeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}

func TestSequencerBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	seq1 := testmodels.NewSequencer(uuid.NewString(), "MiSeq", "Illumina", true)
	seq2 := testmodels.NewSequencer(uuid.NewString(), "NextSeq", "Illumina", true)
	sequencers := []models.Sequencer{seq1, seq2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewSequencerSeedRepository(db)

		err := repo.BulkInsert(context.Background(), sequencers)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewSequencerSeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), sequencers)

		assert.Error(t, err)
	})
}

func TestSequencerCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewSequencerSeedRepository(db)

	seq1 := testmodels.NewSequencer(uuid.NewString(), "MiSeq", "Illumina", true)
	seq2 := testmodels.NewSequencer(uuid.NewString(), "NextSeq", "Illumina", true)
	db.Create(&seq1)
	db.Create(&seq2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewSequencerSeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}

func TestLaboratoryBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	lab1 := testmodels.NewLaboratory(uuid.NewString(), "LACEN/RJ", "RJ", true)
	lab2 := testmodels.NewLaboratory(uuid.NewString(), "LACEN/MG", "MG", true)
	laboratories := []models.Laboratory{lab1, lab2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewLaboratorySeedRepository(db)

		err := repo.BulkInsert(context.Background(), laboratories)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewLaboratorySeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), laboratories)

		assert.Error(t, err)
	})
}

func TestLaboratoryCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewLaboratorySeedRepository(db)

	lab1 := testmodels.NewLaboratory(uuid.NewString(), "LACEN/RJ", "RJ", true)
	lab2 := testmodels.NewLaboratory(uuid.NewString(), "LACEN/MG", "MG", true)
	db.Create(&lab1)
	db.Create(&lab2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewLaboratorySeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}

func TestSampleSourceBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	source1 := testmodels.NewSampleSource(uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Respiratório", "en": "Respiratory", "es": "Respiratorio"},
		true)
	source2 := testmodels.NewSampleSource(uuid.NewString(),
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true)
	sampleSources := []models.SampleSource{source1, source2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewSampleSourceSeedRepository(db)

		err := repo.BulkInsert(context.Background(), sampleSources)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewSampleSourceSeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), sampleSources)

		assert.Error(t, err)
	})
}

func TestSampleSourceCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewSampleSourceSeedRepository(db)

	source1 := testmodels.NewSampleSource(uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Respiratório", "en": "Respiratory", "es": "Respiratorio"},
		true)
	source2 := testmodels.NewSampleSource(uuid.NewString(),
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true)
	db.Create(&source1)
	db.Create(&source2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewSampleSourceSeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}

func TestHealthServiceBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	country := testmodels.NewCountry("BRA", nil)
	db.Create(&country)

	hs1 := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John", "john@example.com", "123456789", true,
	)
	hs2 := testmodels.NewHealthService(
		uuid.NewString(), "Hospital B", models.Private, country,
		"São Paulo", "Jane", "jane@example.com", "987654321", true,
	)
	healthServices := []models.HealthService{hs1, hs2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewHealthServiceSeedRepository(db)

		err := repo.BulkInsert(context.Background(), healthServices)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewHealthServiceSeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), healthServices)

		assert.Error(t, err)
	})
}

func TestHealthServiceCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewHealthServiceSeedRepository(db)

	country := testmodels.NewCountry("BRA", nil)
	db.Create(&country)

	hs1 := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John", "john@example.com", "123456789", true,
	)
	hs2 := testmodels.NewHealthService(
		uuid.NewString(), "Hospital B", models.Private, country,
		"São Paulo", "Jane", "jane@example.com", "987654321", true,
	)
	db.Create(&hs1)
	db.Create(&hs2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewHealthServiceSeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}
