package responses_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
)

func TestGetResponse(t *testing.T) {
	translation.LoadTranslation()

	localizerPt := i18n.NewLocalizer(translation.Bundle, "pt")
	localizerEn := i18n.NewLocalizer(translation.Bundle, "en")
	localizerEs := i18n.NewLocalizer(translation.Bundle, "es")

	expectedPt := "Login efetuado com sucesso."
	expectedEn := "Login successful."
	expectedEs := "Inicio de sesión exitoso."

	expecteds := []string{expectedPt, expectedEn, expectedEs}

	resultPt := responses.GetResponse(localizerPt, responses.LoginSuccess)
	resultEn := responses.GetResponse(localizerEn, responses.LoginSuccess)
	resultEs := responses.GetResponse(localizerEs, responses.LoginSuccess)

	results := []string{resultPt, resultEn, resultEs}

	assert.Equal(t, expecteds, results)
}

func TestGetResponseWithData(t *testing.T) {
	translation.LoadTranslation()

	localizerPt := i18n.NewLocalizer(translation.Bundle, "pt")
	localizerEn := i18n.NewLocalizer(translation.Bundle, "en")
	localizerEs := i18n.NewLocalizer(translation.Bundle, "es")

	expectedPt := "O nome de usuário deve ter pelo menos 5 caracteres."
	expectedEn := "Username must be at least 5 characters long."
	expectedEs := "El nombre de usuario debe tener al menos 5 caracteres."

	expecteds := []string{expectedPt, expectedEn, expectedEs}

	data := map[string]any{"Param": 5}

	resultPt := responses.GetResponseWithData(localizerPt, "validation.Username.min", data)
	resultEn := responses.GetResponseWithData(localizerEn, "validation.Username.min", data)
	resultEs := responses.GetResponseWithData(localizerEs, "validation.Username.min", data)

	results := []string{resultPt, resultEn, resultEs}

	assert.Equal(t, expecteds, results)
}
