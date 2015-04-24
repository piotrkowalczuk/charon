package lib

// ValidationErrorBuilder ...
type ValidationErrorBuilder struct {
	hasErrors bool
	errors    map[string]*ValidationError
}

// NewValidationErrorBuilder ...
func NewValidationErrorBuilder() *ValidationErrorBuilder {
	return &ValidationErrorBuilder{
		hasErrors: false,
		errors:    make(map[string]*ValidationError),
	}
}

// Add ...
func (veb *ValidationErrorBuilder) Add(key, message string) {
	if veb.exists(key, message) {
		return
	}

	if veb.errors[key] == nil {
		veb.errors[key] = &ValidationError{}
	}

	veb.hasErrors = true
	veb.errors[key].add(message)
}

func (veb *ValidationErrorBuilder) exists(key, message string) bool {
	if _, exists := veb.errors[key]; !exists {
		return false
	}

	return veb.errors[key].exists(message)
}

// HasErrors ...
func (veb *ValidationErrorBuilder) HasErrors() bool {
	return veb.hasErrors
}

// Errors ...
func (veb *ValidationErrorBuilder) Errors() map[string]*ValidationError {
	return veb.errors
}
