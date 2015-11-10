package request

import (
	"net/url"

	validator "github.com/asaskevich/govalidator"
	"github.com/piotrkowalczuk/charon/charond/lib"
)

// PasswordRecoveryConfirmationRequest ...
type PasswordRecoveryConfirmationRequest struct {
	Password       string
	PasswordRepeat string
}

// PasswordRecoveryConfirmationRequestFromForm ...
func NewPasswordRecoveryConfirmationRequestFromForm(form url.Values) *PasswordRecoveryConfirmationRequest {
	return &PasswordRecoveryConfirmationRequest{
		Password:       form.Get("password"),
		PasswordRepeat: form.Get("passwordRepeat"),
	}
}

// Validate ...
func (rr *PasswordRecoveryConfirmationRequest) Validate(builder *lib.ValidationErrorBuilder) {
	if !validator.IsByteLength(rr.Password, 6, 45) {
		builder.Add("password", "validation_error.password_min_max_wrong")
	}

	if rr.Password != rr.PasswordRepeat {
		builder.Add("passwordRepeat", "validation_error.password_repeat_mismatch")
	}
}
