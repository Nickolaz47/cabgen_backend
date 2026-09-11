package middlewares

import (
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := translation.ParseLanguage(c.GetHeader("Accept-Language"))
		localizer := i18n.NewLocalizer(translation.Bundle, lang)

		c.Set(translation.LocalizerKey, localizer)
		c.Set("lang", lang)

		c.Next()
	}
}
