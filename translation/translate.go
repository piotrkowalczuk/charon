package translation

import (
	"errors"
	"fmt"
)

// SourceStrategy ...
type SourceStrategy interface {
	Get(lang string, key string) (string, error)
}

// Translate ...
type Translate struct {
	source      SourceStrategy
	defaultLang string
}

// LangDecider ...
type LangDecider func() (string, error)

// TransFunc which can be used in templating
type TransFunc func(string, ...interface{}) (string, error)

// ErrNoTranslation ...
var ErrNoTranslation = errors.New("translation: No Translation for given key")

// Get ...
func (t *Translate) Get(lang string, key string) (string, error) {
	val, err := t.source.Get(lang, key)

	if err != nil || val == "" {
		return key, err
	}

	return val, nil
}

// GetTransFunc returns function that can be used in templating
// first arg of returned function takes translation key, rest works like arguments to fmt.Sprintf
func (t *Translate) GetTransFunc(ld LangDecider) TransFunc {
	return func(toTranslate string, a ...interface{}) (string, error) {
		lang, err := ld()

		if err != nil {
			return "", err
		}

		translated, err := t.Get(lang, toTranslate)

		if err != nil && err != ErrNoTranslation {
			return translated, err
		}

		return fmt.Sprintf(translated, a...), nil
	}
}

// New ...
func New(sourceStrategy SourceStrategy, defaultLang string) *Translate {
	return &Translate{
		source:      sourceStrategy,
		defaultLang: defaultLang,
	}
}
