package i18n

import (
	"sync"

	"github.com/memnix/memnix-rest/pkg/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	// Bundle is the i18n bundle.
	bundle *i18n.Bundle //nolint:gochecknoglobals //Singleton
	once   sync.Once    //nolint:gochecknoglobals //Singleton
)

// GetBundle returns the i18n bundle.
func GetBundle() *i18n.Bundle {
	once.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.NewJSON(&json.SonicJSON{}).Unmarshal)

		// This code should be improved so that it can be tested and follow a configuration file
		bundle.MustLoadMessageFile("locales/en.json")
		bundle.MustLoadMessageFile("locales/fr.json")
	})
	return bundle
}

// Localizer returns a localizer for the given language and accept language.
func Localizer(lang, accept string) *i18n.Localizer {
	return i18n.NewLocalizer(GetBundle(), lang, accept)
}

func DefaultLocalizer() *i18n.Localizer {
	return Localizer("en", "en")
}

func MockLocalizer(path string) *i18n.Localizer {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("json", json.NewJSON(&json.SonicJSON{}).Unmarshal)
	b.MustLoadMessageFile(path)
	return i18n.NewLocalizer(b, "en", "en")
}
