package web

import (
	"net/http"

	"github.com/go-soa/charon/controller/web/request"
	"github.com/go-soa/charon/lib"
	"github.com/go-soa/charon/lib/security"
	"github.com/go-soa/charon/model"
	"github.com/go-soa/charon/repository"
	"golang.org/x/net/context"
)

// RegistrationIndex ...
func (h *Handler) RegistrationIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	err := h.Container.Templates.ExecuteTemplate(rw, h.TemplateName, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RegistrationCreate ...
func (h *Handler) RegistrationCreate(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	validationErrorBuilder := lib.NewValidationErrorBuilder()

	registrationRequest := request.NewRegistrationRequestFromForm(r.Form)
	registrationRequest.Validate(validationErrorBuilder)

	if validationErrorBuilder.HasErrors() {
		rw.WriteHeader(http.StatusBadRequest)
		err := h.Container.Templates.ExecuteTemplate(rw, h.TemplateName, map[string]interface{}{
			"validationErrors": validationErrorBuilder.Errors(),
			"request":          registrationRequest,
		})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	user, err := createAndRegisterUser(h.Container.PasswordHasher, h.Container.RM.User, registrationRequest)
	if err != nil {
		if err == repository.ErrUserUniqueConstraintViolationUsername {
			validationErrorBuilder.Add("email", "User with given email already exists.")

			err = h.Container.Templates.ExecuteTemplate(rw, h.TemplateName, map[string]interface{}{
				"validationErrors": validationErrorBuilder.Errors(),
				"request":          registrationRequest,
			})
		}

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	err = h.Container.Mailer.SendWelcomeMail(user.Username, user.String())

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/registration/success", http.StatusFound)
}

// RegistrationSuccess ...
func (h *Handler) RegistrationSuccess(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	err := h.Container.Templates.ExecuteTemplate(rw, h.TemplateName, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createAndRegisterUser(
	passwordHasher security.PasswordHasher,
	repository *repository.UserRepository,
	request *request.RegistrationRequest,
) (*model.User, error) {
	hashedPassword, err := passwordHasher.Hash(request.Password)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(request.Email, hashedPassword, request.FirstName, request.LastName)

	_, err = repository.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
