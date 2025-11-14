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

	updateInput := models.AdminUpdateInput{
		UpdateUserInput: models.UpdateUserInput{
			Name:        &name,
			Username:    &username,
			Institution: &institution,
			Interest:    &interest,
			Role:        &role,
		},
		Email: &email,
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
