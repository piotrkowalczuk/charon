package request

import (
	"net/url"

	validator "github.com/asaskevich/govalidator"
	"github.com/go-soa/auth/lib"
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
		builder.Add("email", "Email address should contain minimum 6 and maximum 45 characters.")
	} else if !validator.IsEmail(rr.Email) {
		builder.Add("email", "Invalid email address.")
	}

	if !validator.IsByteLength(rr.Password, 6, 45) {
		builder.Add("password", "Password should contain minimum 6 and maximum 45 characters .")
	}

	if !validator.IsByteLength(rr.FirstName, 1, 45) {
		builder.Add("firstName", "First name should contain maximum 45 characters .")
	}

	if !validator.IsByteLength(rr.LastName, 1, 45) {
		builder.Add("lastName", "Last name should contain maximum 45 characters .")
	}
}
