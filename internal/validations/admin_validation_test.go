package validations_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestApplyAdminUpdateToUser(t *testing.T) {
	user := models.User{}

	name := "Nicolas Silva"
	username := "nikol"
	email := "nicolas@mail.com"
	institution := "Fiocruz"
	interest := "Programming"
	role := "Developer"

	updateInput := models.AdminUserUpdateInput{
		Name:        &name,
		Username:    &username,
		Institution: &institution,
		Interest:    &interest,
		Role:        &role,
		Email:       &email,
	}

	validations.ApplyAdminUpdateToUser(&user, &updateInput)

	assert.Equal(t, name, user.Name)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, &institution, user.Institution)
	assert.Equal(t, &interest, user.Interest)
	assert.Equal(t, &role, user.Role)
}

func TestValidateTranslationMap(t *testing.T) {
	testutils.SetupTestContext()
	c, _ := testutils.SetupGinContext(
		http.MethodGet, "/", "",
		nil, nil,
	)

	t.Run("Success - Origin", func(t *testing.T) {
		names := map[string]string{
			"pt": "Humano",
			"en": "Human",
			"es": "Humano",
		}

		errMsg, ok := validations.ValidateTranslationMap(c, "origin", names)

		assert.Empty(t, errMsg)
		assert.True(t, ok)
	})

	t.Run("Success - Sample Source", func(t *testing.T) {
		names := map[string]string{
			"pt": "Sangue",
			"en": "Blood",
			"es": "Sangre",
		}

		errMsg, ok := validations.ValidateTranslationMap(c, "sampleSource", names)

		assert.Empty(t, errMsg)
		assert.True(t, ok)
	})

	t.Run("Missing language", func(t *testing.T) {
		names := map[string]string{
			"pt": "Humano",
			"en": "Human",
		}

		errMsg, ok := validations.ValidateTranslationMap(c, "origin", names)

		assert.Equal(t, errMsg, "Missing es translation.")
		assert.False(t, ok)
	})

	t.Run("Empty translation", func(t *testing.T) {
		names := map[string]string{
			"pt": "",
			"en": "Human",
			"es": "Humano",
		}

		errMsg, ok := validations.ValidateTranslationMap(c, "origin", names)

		assert.Equal(t, errMsg, "Empty pt translation.")
		assert.False(t, ok)
	})
}

func TestApplyOriginUpdate(t *testing.T) {
	origin := models.Origin{
		ID:       uuid.New(),
		Names:    map[string]string{"pt": "Humano", "en": "Human", "es": "Human"},
		IsActive: false,
	}

	isActive := true
	originUpdate := models.OriginUpdateInput{
		Names:    map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		IsActive: &isActive,
	}

	expected := models.Origin{
		ID:       origin.ID,
		Names:    originUpdate.Names,
		IsActive: *originUpdate.IsActive,
	}

	validations.ApplyOriginUpdate(&origin, &originUpdate)

	assert.Equal(t, expected, origin)
}

func TestApplySequencerUpdate(t *testing.T) {
	sequencer := models.Sequencer{
		ID:       uuid.New(),
		Brand:    "Ilumina",
		Model:    "Myseq",
		IsActive: false,
	}

	brand := "Illumina"
	model := "MiSeq"
	isActive := true
	sequencerUpdate := models.SequencerUpdateInput{
		Brand:    &brand,
		Model:    &model,
		IsActive: &isActive,
	}

	expected := models.Sequencer{
		ID:       sequencer.ID,
		Brand:    *sequencerUpdate.Brand,
		Model:    *sequencerUpdate.Model,
		IsActive: *sequencerUpdate.IsActive,
	}

	validations.ApplySequencerUpdate(&sequencer, &sequencerUpdate)

	assert.Equal(t, expected, sequencer)
}

func TestApplySampleSourceUpdate(t *testing.T) {
	sampleSource := models.SampleSource{
		ID: uuid.New(),
		Names: map[string]string{
			"pt": "Plasma",
			"en": "Plasm",
			"es": "Plasme",
		},
		Groups: map[string]string{
			"pt": "Sangue",
			"en": "Blood",
			"es": "Sangre",
		},
		IsActive: false,
	}

	isActive := true
	sampleSourceUpdate := models.SampleSourceUpdateInput{
		Names: map[string]string{
			"pt": "Plasma",
			"en": "Plasma",
			"es": "Plasma",
		},
		IsActive: &isActive,
	}

	expected := models.SampleSource{
		ID:       sampleSource.ID,
		Names:    sampleSourceUpdate.Names,
		Groups:   sampleSource.Groups,
		IsActive: *sampleSourceUpdate.IsActive,
	}

	validations.ApplySampleSourceUpdate(&sampleSource, &sampleSourceUpdate)

	assert.Equal(t, expected, sampleSource)
}

func TestApplyLaboratoryUpdate(t *testing.T) {
	laboratory := testmodels.NewLaboratory(
		uuid.NewString(),
		"Laboratório do RJ",
		"LABJ",
		false,
	)

	name := "Laboratório do Rio de Janeiro"
	abbreviation := "LABRJ"
	isActive := true
	laboratoryUpdate := models.LaboratoryUpdateInput{
		Name:         &name,
		Abbreviation: &abbreviation,
		IsActive:     &isActive,
	}

	expected := models.Laboratory{
		ID:           laboratory.ID,
		Name:         *laboratoryUpdate.Name,
		Abbreviation: *laboratoryUpdate.Abbreviation,
		IsActive:     *laboratoryUpdate.IsActive,
	}

	validations.ApplyLaboratoryUpdate(&laboratory, &laboratoryUpdate)

	assert.Equal(t, expected, laboratory)
}

func TestApplyCountryUpdate(t *testing.T) {
	country := testmodels.NewCountry("", nil)

	code, names := "SPN", map[string]string{
		"pt": "Espanha",
		"en": "Spain",
		"es": "España",
	}
	input := models.CountryUpdateInput{
		Code:  &code,
		Names: names,
	}

	expected := models.Country{
		Code:  code,
		Names: names,
	}

	validations.ApplyCountryUpdate(&country, &input)

	assert.Equal(t, expected, country)
}

func TestApplyMicroorganismUpdate(t *testing.T) {
	microorganism := testmodels.NewMicroorganism(
		uuid.NewString(), models.Virus, "Flavivirus",
		nil, false)

	taxon := models.Bacteria
	species := "E. coli"
	variety := map[string]string{
		"pt": "Variedade",
		"en": "Variety",
		"es": "Variedad",
	}
	isActive := true

	input := models.MicroorganismUpdateInput{
		Taxon:    &taxon,
		Species:  &species,
		Variety:  variety,
		IsActive: &isActive,
	}

	expected := models.Microorganism{
		ID:       microorganism.ID,
		Taxon:    *input.Taxon,
		Species:  *input.Species,
		Variety:  input.Variety,
		IsActive: *input.IsActive,
	}

	validations.ApplyMicroorganismUpdate(&microorganism, &input)

	assert.Equal(t, expected, microorganism)
}

func TestApplyHealthServiceUpdate(t *testing.T) {
	country := testmodels.NewCountry("", nil)

	healthService := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John Doe", "john@example.com", "123456789",
		false,
	)

	name := "Hospital B"
	typeStr := models.Private
	city := "Sao Paulo"
	contactant := "Jane Doe"
	contactEmail := "jane@example.com"
	contactPhone := "987654321"
	isActive := true

	input := models.HealthServiceUpdateInput{
		Name:         &name,
		Type:         &typeStr,
		City:         &city,
		Contactant:   &contactant,
		ContactEmail: &contactEmail,
		ContactPhone: &contactPhone,
		IsActive:     &isActive,
	}

	expected := models.HealthService{
		ID:           healthService.ID,
		Name:         *input.Name,
		Type:         *input.Type,
		CountryID:    country.ID,
		Country:      country,
		City:         *input.City,
		Contactant:   *input.Contactant,
		ContactEmail: *input.ContactEmail,
		ContactPhone: *input.ContactPhone,
		IsActive:     *input.IsActive,
	}

	validations.ApplyHealthServiceUpdate(&healthService, &input)

	assert.Equal(t, expected, healthService)
}
