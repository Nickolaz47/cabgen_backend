package middlewares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
)

func TestI18nMiddleware(t *testing.T) {
	testutils.SetupTestContext()

	tests := []struct {
		name    string
		header  string
		lang    string
		message string
	}{
		{"pt", "pt", "pt", "Usuário não encontrado."},
		{"en", "en", "en", "User not found."},
		{"es", "es", "es", "Usuario no encontrado."},
		{"variant pt-BR", "pt-BR", "pt", "Usuário não encontrado."},
		{"uppercase PT", "PT", "pt", "Usuário não encontrado."},
		{"variant with q", "pt-BR;q=0.9", "pt", "Usuário não encontrado."},
		{"list with q", "pt-BR,es;q=0.9,en;q=0.8", "pt", "Usuário não encontrado."},
		{"en-US variant", "en-US", "en", "User not found."},
		{"unsupported falls back to en", "fr", "en", "User not found."},
		{"empty header falls back to en", "", "en", "User not found."},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Success %s", tt.name), func(t *testing.T) {
			w, r := testutils.SetupMiddlewareContext()

			testutils.AddMiddlewares(r, middlewares.I18nMiddleware())

			r.GET("/", func(c *gin.Context) {
				v, exists := c.Get(translation.LocalizerKey)
				localizer, ok := v.(*i18n.Localizer)

				rawLanguage, langExists := c.Get("lang")
				language, langOK := rawLanguage.(string)

				if exists && ok && langExists && langOK {
					c.JSON(http.StatusOK, map[string]any{
						"translatedMessage": responses.GetResponse(
							localizer, responses.UserNotFoundError),
						"language": language,
					})
					return
				}
				c.Status(http.StatusNotFound)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Accept-Language", tt.header)
			}
			r.ServeHTTP(w, req)

			expected := fmt.Sprintf(
				`{"translatedMessage": "%s", "language": "%s"}`,
				tt.message, tt.lang)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, expected, w.Body.String())
		})
	}
}
