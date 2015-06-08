package translation

import "errors"

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
type TransFunc func(string) (string, error)

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
func (t *Translate) GetTransFunc(ld LangDecider) TransFunc {
	return func(toTranslate string) (string, error) {
		lang, err := ld()

		if err != nil {
			return "", err
		}

		translated, err := t.Get(lang, toTranslate)

		if err != nil && err != ErrNoTranslation {
			return translated, err
		}

		return translated, nil
	}
}

// New ...
func New(sourceStrategy SourceStrategy, defaultLang string) *Translate {
	return &Translate{
		source:      sourceStrategy,
		defaultLang: defaultLang,
	}
}
