package web

import (
	"net/http"

	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/piotrkowalczuk/charon/charond/controller/web/request"
	"github.com/piotrkowalczuk/charon/charond/lib"
	mnemosyne "github.com/piotrkowalczuk/mnemosyne/lib"
	"golang.org/x/net/context"
)

// LoginIndex ...
func (h *Handler) LoginIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	return h.renderTemplate(rw, ctx)
}

// LoginProcess ...
func (h *Handler) LoginProcess(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	r.ParseForm()

	validationErrorBuilder := lib.NewValidationErrorBuilder()

	loginRequest := request.NewLoginRequest(r.Form)
	loginRequest.Validate(validationErrorBuilder)

	if validationErrorBuilder.HasErrors() {
		rw.WriteHeader(http.StatusBadRequest)
		return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
			"validationErrors": validationErrorBuilder.Errors(),
			"request":          loginRequest,
		})
	}

	user, err := h.Container.RM.User.FindOneByUsername(loginRequest.Email)
	if err != nil {
		return h.renderTemplate500(rw, ctx, err)
	}

	if matches := h.Container.PasswordHasher.Compare(user.Password, loginRequest.Password); !matches {
		h.Container.Logger.WithFields(logrus.Fields{
			"username": loginRequest.Email,
			"password": loginRequest.Password,
		}).Debug("Wrong password provided.")

		validationErrorBuilder.Add("form", "The email address and password you entered do not match.")
		return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
			"validationErrors": validationErrorBuilder.Errors(),
			"request":          loginRequest,
		})
	}

	if !user.IsConfirmed || !user.IsActive {
		return h.renderTemplate403(rw, ctx)
	}

	session := mnemosyne.Session{}
	sessionData := mnemosyne.SessionData{
		"user_id":    strconv.FormatInt(user.ID, 10),
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}

	err = h.Container.Mnemosyne.Call("Store.New", sessionData, &session)
	if err != nil {
		return h.renderTemplate500(rw, ctx, err)
	}

	err = h.Container.RM.User.UpdateLastLoginAt(user.ID)
	if err != nil {
		return h.renderTemplate500(rw, ctx, err)
	}

	cookie := &http.Cookie{
		Name:     "sid",
		Value:    session.ID.String(),
		HttpOnly: true,
		// TODO: need to be checked what behavior is preferred
		//		Domain:   h.Container.Config.Domain,
	}

	http.SetCookie(rw, cookie)

	h.redirect(rw, r, "dashboard", http.StatusFound)

	return ctx
}

// LogoutIndex ...
func (h *Handler) LogoutIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	cookie, err := r.Cookie("sid")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return h.renderTemplate400(rw, ctx)
		default:
			return h.renderTemplate500(rw, ctx, err)
		}
	}

	err = h.Container.Mnemosyne.Call("Store.Abandon", mnemosyne.SessionID(cookie.Value), nil)
	if err != nil {
		switch err {
		case mnemosyne.ErrSessionNotFound:
			return h.renderTemplate403(rw, ctx)
		default:
			return h.renderTemplate500(rw, ctx, err)
		}
	}

	h.redirect(rw, r, "login", http.StatusFound)

	return ctx
}
