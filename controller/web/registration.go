package web

import (
	"net/http"

	"github.com/go-soa/auth/controller/web/request"
	"github.com/go-soa/auth/model"
	"github.com/go-soa/auth/repository"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// RegistrationIndex ...
func (h *Handler) RegistrationIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(rw, h.TmplName, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RegistrationCreate ...
func (h *Handler) RegistrationCreate(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	registrationRequest := request.NewRegistrationRequestFromForm(r.Form)
	validationErrorBuilder := registrationRequest.Validate()

	if validationErrorBuilder.HasErrors() {
		rw.WriteHeader(http.StatusBadRequest)
		err := h.Tmpl.ExecuteTemplate(rw, "registration_index", map[string]interface{}{
			"validationErrors": validationErrorBuilder.Errors(),
			"request":          registrationRequest,
		})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	user, err := createAndRegisterUser(h.RM.User, registrationRequest)
	if err != nil {
		if err == repository.ErrUserUniqueConstraintViolationUsername {
			validationErrorBuilder.Add("email", "User with given email already exists.")

			rw.WriteHeader(http.StatusBadRequest)
			err = h.Tmpl.ExecuteTemplate(rw, "registration_index", map[string]interface{}{
				"validationErrors": validationErrorBuilder.Errors(),
				"request":          registrationRequest,
			})
		}

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	err = h.Mailer.SendWelcomeMail(user.Username, user.String())

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/registration/success", http.StatusMovedPermanently)
}

// RegistrationSuccess ...
func (h *Handler) RegistrationSuccess(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(rw, h.TmplName, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createAndRegisterUser(
	repository *repository.UserRepository,
	request *request.RegistrationRequest,
) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(request.Email, string(hashedPassword), request.FirstName, request.LastName)

	_, err = repository.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
