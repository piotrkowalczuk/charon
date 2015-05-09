package web

import (
	"errors"
	"net/http"

	"strconv"

	"github.com/go-soa/charon/controller/web/request"
	"github.com/go-soa/charon/lib"
	"github.com/go-soa/charon/lib/routing"
	"github.com/go-soa/charon/lib/security"
	"github.com/go-soa/charon/model"
	"github.com/go-soa/charon/repository"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// RegistrationIndex ...
func (h *Handler) RegistrationIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	h.renderTemplate(rw)
}

// RegistrationCreate ...
func (h *Handler) RegistrationCreate(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	validationErrorBuilder := lib.NewValidationErrorBuilder()

	registrationRequest := request.NewRegistrationRequestFromForm(r.Form)
	registrationRequest.Validate(validationErrorBuilder)

	if validationErrorBuilder.HasErrors() {
		rw.WriteHeader(http.StatusBadRequest)
		h.renderTemplateWithData(rw, map[string]interface{}{
			"validationErrors": validationErrorBuilder.Errors(),
			"request":          registrationRequest,
		})

		return
	}

	user, err := createAndRegisterUser(h.Container.PasswordHasher, h.Container.RM.User, registrationRequest)
	if err != nil {
		switch err {
		case repository.ErrUserUniqueConstraintViolationUsername:
			validationErrorBuilder.Add("email", "User with given email already exists.")

			h.renderTemplateWithData(rw, map[string]interface{}{
				"validationErrors": validationErrorBuilder.Errors(),
				"request":          registrationRequest,
			})
		default:
			h.sendError500(rw, err)
		}

		return
	}

	err = h.Container.ConfirmationMailer.Send(user.Username, map[string]interface{}{
		"user": user,
	})

	if err != nil {
		h.sendError500(rw, err)
		return
	}

	http.Redirect(rw, r, "/registration/success", http.StatusFound)
}

// RegistrationSuccess ...
func (h *Handler) RegistrationSuccess(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	h.renderTemplate(rw)
}

// RegistrationConfirmation ...
func (h *Handler) RegistrationConfirmation(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	var ok bool
	var confirmationTokenParam string
	var userIdParam string

	if confirmationTokenParam, ok = routing.ParamFromContext(ctx, "confirmationToken"); !ok {
		h.sendError400(rw, errors.New("controller/web: confirmationToken param is missing"))
	}

	if userIdParam, ok = routing.ParamFromContext(ctx, "userId"); !ok {
		h.sendError400(rw, errors.New("controller/web: userId param is missing"))
	}

	userID, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		h.sendError400(rw, errors.New("controller/web: userId wrong type"))
	}

	if err := h.Container.RM.User.RegistrationConfirmation(userID, confirmationTokenParam); err != nil {
		switch err {
		case repository.ErrUserNotFound:
			h.renderTemplateWithStatus(rw, http.StatusMethodNotAllowed)
		default:
			h.sendError500(rw, err)
		}

		return
	}

	h.renderTemplate(rw)
}

func createAndRegisterUser(
	passwordHasher security.PasswordHasher,
	repository *repository.UserRepository,
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

	_, err = repository.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
