package translation

// SourceStrategy ...
type SourceStrategy interface {
	Get(lang string, key string) (string, error)
}

// Translate ...
type Translate struct {
	source      SourceStrategy
	defaultLang string
}

// Get ...
func (t *Translate) Get(lang string, key string) (string, error) {
	val, err := t.source.Get(lang, key)

	if err != nil || val == "" {
		return key, err
	}

	return val, nil
}

// New ...
func New(sourceStrategy SourceStrategy, defaultLang string) *Translate {
	return &Translate{
		source:      sourceStrategy,
		defaultLang: defaultLang,
	}
}
