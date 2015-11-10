package lib

// Validator ...
type Validator interface {
	Validate(*ValidationErrorBuilder)
}
