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
