package validations_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
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

func TestGetOffset(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success - Absent defaults to zero", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)

		offset, ok := validations.GetOffset(c)

		assert.True(t, ok)
		assert.Equal(t, 0, offset)
	})

	t.Run("Success - Valid offset", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/?offset=250",
			"", nil, nil)

		offset, ok := validations.GetOffset(c)

		assert.True(t, ok)
		assert.Equal(t, 250, offset)
	})

	t.Run("Error - Non numeric", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/?offset=abc",
			"", nil, nil)

		offset, ok := validations.GetOffset(c)

		assert.False(t, ok)
		assert.Equal(t, 0, offset)
	})

	t.Run("Error - Negative", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/?offset=-1",
			"", nil, nil)

		offset, ok := validations.GetOffset(c)

		assert.False(t, ok)
		assert.Equal(t, 0, offset)
	})
}
