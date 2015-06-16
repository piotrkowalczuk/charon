package request

import (
	"net/url"

	validator "github.com/asaskevich/govalidator"
	"github.com/go-soa/charon/lib"
)

// LoginRequest ...
type LoginRequest struct {
	Email    string
	Password string
}

// NewLoginRequest ...
func NewLoginRequest(form url.Values) *LoginRequest {
	return &LoginRequest{
		Email:    form.Get("email"),
		Password: form.Get("password"),
	}
}

// Validate ...
func (lr *LoginRequest) Validate(builder *lib.ValidationErrorBuilder) {
	if !validator.IsByteLength(lr.Email, 1, 255) {
		builder.Add("email", "Email address is missing.")
	}

	if !validator.IsByteLength(lr.Password, 1, 255) {
		builder.Add("password", "Password is missing.")
	}
}
