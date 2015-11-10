package request

import (
	"net/url"

	validator "github.com/asaskevich/govalidator"
	"github.com/piotrkowalczuk/charon/charond/lib"
)

// PasswordRecoveryRequest ...
type PasswordRecoveryRequest struct {
	Email string
}

// NewPasswordRecoveryRequestFromForm ...
func NewPasswordRecoveryRequestFromForm(form url.Values) *PasswordRecoveryRequest {
	return &PasswordRecoveryRequest{
		Email: form.Get("email"),
	}
}

// Validate ...
// TODO: unify somehow, to have centralized place where email (or any other field) is validated
func (prr *PasswordRecoveryRequest) Validate(builder *lib.ValidationErrorBuilder) {
	if !validator.IsByteLength(prr.Email, 6, 45) {
		builder.Add("email", "validation_error.email_min_max_wrong")
	} else if !validator.IsEmail(prr.Email) {
		builder.Add("email", "Invalid email address.")
	}
}
