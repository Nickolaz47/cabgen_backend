package translation_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
)

func TestLoadTranslation(t *testing.T) {
	assert.Empty(t, translation.Bundle)

	translation.LoadTranslation()

	assert.NotEmpty(t, translation.Bundle)
}

func TestGetLocalizerFromContext(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Language header not empty", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(
			http.MethodGet, "/", "",
			nil, nil,
		)
		expectedLocalizer := i18n.NewLocalizer(translation.Bundle, "pt")
		c.Set(translation.LocalizerKey, expectedLocalizer)

		localizer := translation.GetLocalizerFromContext(c)

		assert.Equal(t, expectedLocalizer, localizer)
	})

	t.Run("Language header empty", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)

		localizer := translation.GetLocalizerFromContext(c)
		expectedLocalizer := i18n.NewLocalizer(translation.Bundle, "en")

		assert.Equal(t, expectedLocalizer, localizer)
	})
}
