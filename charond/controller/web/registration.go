package web

import (
	"net/http"

	"strconv"

	"github.com/piotrkowalczuk/charon/charond/controller/web/request"
	"github.com/piotrkowalczuk/charon/charond/lib"
	"github.com/piotrkowalczuk/charon/charond/lib/routing"
	"github.com/piotrkowalczuk/charon/charond/lib/security"
	"github.com/piotrkowalczuk/charon/charond/model"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// RegistrationIndex ...
func (h *Handler) RegistrationIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	return h.renderTemplate(rw, ctx)
}

// RegistrationProcess ...
func (h *Handler) RegistrationProcess(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	r.ParseForm()

	validationErrorBuilder := lib.NewValidationErrorBuilder()

	registrationRequest := request.NewRegistrationRequestFromForm(r.Form)
	registrationRequest.Validate(validationErrorBuilder)

	if validationErrorBuilder.HasErrors() {
		rw.WriteHeader(http.StatusBadRequest)
		return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
			"validationErrors": validationErrorBuilder.Errors(),
			"request":          registrationRequest,
		})
	}

	user, err := createAndRegisterUser(h.Container.PasswordHasher, h.Container.RM.User, registrationRequest)
	if err != nil {
		switch err {
		case lib.ErrUserUniqueConstraintViolationUsername:
			validationErrorBuilder.Add("email", "User with given email already exists.")

			return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
				"validationErrors": validationErrorBuilder.Errors(),
				"request":          registrationRequest,
			})
		default:
			return h.renderTemplate500(rw, ctx, err)
		}
	}

	err = h.Container.ConfirmationMailer.Send(user.Username, map[string]interface{}{
		"user": user,
	})
	if err != nil {
		return h.renderTemplate500(rw, ctx, err)
	}

	h.redirect(rw, r, "registration_success", http.StatusFound)

	return ctx
}

// RegistrationSuccess ...
func (h *Handler) RegistrationSuccess(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	return h.renderTemplate(rw, ctx)
}

// RegistrationConfirmation ...
func (h *Handler) RegistrationConfirmation(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	var ok bool
	var confirmationTokenParam string
	var userIDParam string

	if confirmationTokenParam, ok = routing.ParamFromContext(ctx, "confirmationToken"); !ok {
		h.Container.Logger.Debug("Confirmation token param is missing.")
		return h.renderTemplate400(rw, ctx)
	}

	if userIDParam, ok = routing.ParamFromContext(ctx, "userId"); !ok {
		h.Container.Logger.Debug("User ID param is missing.")
		return h.renderTemplate400(rw, ctx)
	}

	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		h.Container.Logger.Debug("User ID param wrong type.")
		return h.renderTemplate400(rw, ctx)
	}

	if err := h.Container.RM.User.RegistrationConfirmation(userID, confirmationTokenParam); err != nil {
		switch err {
		case lib.ErrUserNotFound:
			h.Container.Logger.Debug("Registration confirmation failure, user not found.")
			return h.renderTemplateWithStatus(rw, ctx, http.StatusMethodNotAllowed)
		default:
			return h.renderTemplate500(rw, ctx, err)
		}
	}

	return h.renderTemplate(rw, ctx)
}

func createAndRegisterUser(
	passwordHasher security.PasswordHasher,
	repository lib.UserRepository,
	request *request.RegistrationRequest,
) (*model.User, error) {
	confirmationToken := uuid.NewV4().String()
	hashedPassword, err := passwordHasher.Hash(request.Password)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(
		request.Email,
		hashedPassword,
		request.FirstName,
		request.LastName,
		confirmationToken,
	)

	err = repository.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
