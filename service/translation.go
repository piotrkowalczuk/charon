package service

import "github.com/go-soa/charon/translation"

type translationConfig struct {
	DefaultLang string `xml:"default-lang"`
	Strategy    string `xml:"strategy"`
}

const translaterInputStrategyCSV = "csv"

// Translater ...
var Translate *translation.Translate

// InitTranslation ...
func InitTranslation(config translationConfig) {
	var sourceStrategy translation.SourceStrategy

	switch config.Strategy {
	case translaterInputStrategyCSV:
		sourceStrategy = translation.NewCSVSource("data/i18n")
	}

	Translate = translation.New(sourceStrategy, config.DefaultLang)
}
