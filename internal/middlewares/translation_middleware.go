package middlewares

import (
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "en"
		}

		tag := strings.Split(lang, ",")[0]
		localizer := i18n.NewLocalizer(translation.Bundle, tag)

		c.Set(translation.LocalizerKey, localizer)
		c.Set("lang", lang)

		c.Next()
	}
}
