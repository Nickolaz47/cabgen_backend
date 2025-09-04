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

	tests := map[string]string{
		"pt": "Usuário não encontrado.",
		"en": "User not found.",
		"es": "Usuario no encontrado.",
	}

	for lang, message := range tests {
		t.Run(fmt.Sprintf("Success %s", lang), func(t *testing.T) {
			w, r := testutils.SetupMiddlewareContext()

			testutils.AddMiddlewares(r, middlewares.I18nMiddleware())

			r.GET("/", func(c *gin.Context) {
				v, exists := c.Get(translation.LocalizerKey)
				localizer := v.(*i18n.Localizer)
				message := responses.GetResponse(localizer, responses.UserNotFoundError)

				if exists {
					c.JSON(http.StatusOK, map[string]any{"lang": message})
				}
				c.Status(http.StatusNotFound)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept-Language", lang)
			r.ServeHTTP(w, req)

			expected := fmt.Sprintf(`{"lang": "%s"}`, message)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, expected, w.Body.String())
		})
	}
}
