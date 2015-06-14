package service

import "github.com/go-soa/charon/translation"

const (
	translatorInputStrategyCSV = "csv"
)

var (
	// Translate  ...
	Translate *translation.Translate
)

type translationConfig struct {
	DefaultLang string `xml:"default-lang"`
	Strategy    string `xml:"strategy"`
}

// InitTranslation ...
func InitTranslation(config translationConfig) {
	var sourceStrategy translation.SourceStrategy

	switch config.Strategy {
	case translatorInputStrategyCSV:
		sourceStrategy = translation.NewCSVSource("data/i18n")
	}

	Translate = translation.New(sourceStrategy, config.DefaultLang)
}
