package lib

// ValidationError ...
type ValidationError struct {
	hasErrors bool
	errors    []string
}

// HasErrors ...
func (ve *ValidationError) HasErrors() bool {
	return ve.hasErrors
}

// Errors ...
func (ve *ValidationError) Errors() []string {
	return ve.errors
}

func (ve *ValidationError) add(message string) {
	ve.hasErrors = true
	ve.errors = append(ve.errors, message)
}

func (ve *ValidationError) exists(message string) bool {
	for _, err := range ve.errors {
		if err == message {
			return true
		}
	}

	return false
}
