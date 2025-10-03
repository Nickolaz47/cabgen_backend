package translation

import (
	"embed"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

const LocalizerKey = "localizer"

var Languages = []string{"pt", "en", "es"}

//go:embed active/*.toml
var localeFS embed.FS

var Bundle *i18n.Bundle

func LoadTranslation() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	var files []string
	for _, lang := range Languages {
		files = append(files, "active/active."+lang+".toml")
	}

	for _, file := range files {
		if _, err := Bundle.LoadMessageFileFS(localeFS, file); err != nil {
			log.Fatalf("failed to load translation file %s: %v", file, err)
		}
	}
}

func GetLocalizerFromContext(c *gin.Context) *i18n.Localizer {
	value, exists := c.Get(LocalizerKey)
	if !exists {
		return i18n.NewLocalizer(Bundle, "en")
	}

	localizer, ok := value.(*i18n.Localizer)
	if !ok {
		return i18n.NewLocalizer(Bundle, "en")
	}

	return localizer
}

func GetLanguageFromContext(c *gin.Context) string {
	value, exists := c.Get("lang")
	if !exists {
		return "en"
	}

	language, ok := value.(string)
	if !ok {
		return "en"
	}

	return language
}
