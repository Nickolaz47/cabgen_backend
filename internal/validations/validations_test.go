package validations_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		body := testutils.ToJSON(models.LoginInput{
			Username: "nick",
			Password: "12345678",
		})

		c, _ := testutils.SetupGinContext(http.MethodPost, "/", body, nil, nil)
		localizer := i18n.NewLocalizer(translation.Bundle, "pt")

		var input models.LoginInput
		msg, ok := validations.Validate(c, localizer, &input)

		assert.True(t, ok)
		assert.Empty(t, msg)
	})

	t.Run("Error", func(t *testing.T) {
		body := testutils.ToJSON(models.LoginInput{
			Username: "nick",
		})

		c, _ := testutils.SetupGinContext(http.MethodPost, "/", body, nil, nil)
		localizer := i18n.NewLocalizer(translation.Bundle, "pt")

		var input models.LoginInput
		msg, ok := validations.Validate(c, localizer, &input)

		assert.False(t, ok)
		assert.NotEmpty(t, msg)
		assert.Equal(t, "A senha é obrigatória.", msg)
	})
}

func TestGetPage(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success - Absent defaults to one", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)

		page, ok := validations.GetPage(c)

		assert.True(t, ok)
		assert.Equal(t, 1, page)
	})

	t.Run("Success - Valid page", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/?page=3",
			"", nil, nil)

		page, ok := validations.GetPage(c)

		assert.True(t, ok)
		assert.Equal(t, 3, page)
	})

	t.Run("Error - Non numeric", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/?page=abc",
			"", nil, nil)

		page, ok := validations.GetPage(c)

		assert.False(t, ok)
		assert.Equal(t, 0, page)
	})

	t.Run("Error - Zero", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/?page=0",
			"", nil, nil)

		page, ok := validations.GetPage(c)

		assert.False(t, ok)
		assert.Equal(t, 0, page)
	})

	t.Run("Error - Negative", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/?page=-1",
			"", nil, nil)

		page, ok := validations.GetPage(c)

		assert.False(t, ok)
		assert.Equal(t, 0, page)
	})

	t.Run("Error - Above max", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet,
			"/?page=999999999999", "", nil, nil)

		page, ok := validations.GetPage(c)

		assert.False(t, ok)
		assert.Equal(t, 0, page)
	})
}

func TestValidateStruct(t *testing.T) {
	testutils.SetupTestContext()
	localizer := i18n.NewLocalizer(translation.Bundle, "pt")

	inNetwork := true

	validInput := func() models.SampleCreateInput {
		return models.SampleCreateInput{
			OriginCode:      "  BR-RJ-01  ",
			CollectionDate:  models.Date{Time: time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)},
			RunNumber:       "  RUN-2026-XYZ  ",
			RunDate:         models.Date{Time: time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)},
			City:            "  Marica  ",
			InNetwork:       &inNetwork,
			CountryCode:     "BRA",
			OriginID:        uuid.New(),
			SampleSourceID:  uuid.New(),
			MicroorganismID: uuid.New(),
			SequencerID:     uuid.New(),
			LaboratoryID:    uuid.New(),
			HealthServiceID: uuid.New(),
		}
	}

	t.Run("Success - trims input", func(t *testing.T) {
		input := validInput()

		msg, ok := validations.ValidateStruct(localizer, &input)

		assert.True(t, ok)
		assert.Empty(t, msg)
		assert.Equal(t, "BR-RJ-01", input.OriginCode)
		assert.Equal(t, "RUN-2026-XYZ", input.RunNumber)
		assert.Equal(t, "Marica", input.City)
		assert.Equal(t, "BRA", input.CountryCode)
	})

	t.Run("Error - binding tag", func(t *testing.T) {
		input := validInput()
		input.OriginCode = "AB"

		msg, ok := validations.ValidateStruct(localizer, &input)

		assert.False(t, ok)
		assert.NotEmpty(t, msg)
	})

	t.Run("Error - missing collection date", func(t *testing.T) {
		input := validInput()
		input.CollectionDate = models.Date{}

		msg, ok := validations.ValidateStruct(localizer, &input)

		assert.False(t, ok)
		assert.NotEmpty(t, msg)
	})
}
