package request

import (
	"net/url"

	validator "github.com/asaskevich/govalidator"
	"github.com/go-soa/charon/lib"
)

// RegistrationRequest ...
type RegistrationRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// NewRegistrationRequestFromForm ...
func NewRegistrationRequestFromForm(form url.Values) *RegistrationRequest {
	return &RegistrationRequest{
		Email:     form.Get("email"),
		Password:  form.Get("password"),
		FirstName: form.Get("firstName"),
		LastName:  form.Get("lastName"),
	}
}

// Validate ...
func (rr *RegistrationRequest) Validate(builder *lib.ValidationErrorBuilder) {
	if !validator.IsByteLength(rr.Email, 6, 45) {
		builder.Add("email", "validation_error.email_min_max_wrong")
	} else if !validator.IsEmail(rr.Email) {
		builder.Add("email", "Invalid email address.")
	}

	if !validator.IsByteLength(rr.Password, 6, 45) {
		builder.Add("password", "validation_error.password_min_max_wrong")
	}

	if !validator.IsByteLength(rr.FirstName, 1, 45) {
		builder.Add("firstName", "validation_error.firstname_max_len_wrong")
	}

	if !validator.IsByteLength(rr.LastName, 1, 45) {
		builder.Add("lastName", "validation_error.lastname_max_len_wrong")
	}
}
