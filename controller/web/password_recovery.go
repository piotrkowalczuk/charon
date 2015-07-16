package web

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/go-soa/charon/controller/web/request"
	"github.com/go-soa/charon/lib"
	"github.com/go-soa/charon/lib/routing"
	"golang.org/x/net/context"
)

// PasswordRecoveryIndex ...
func (h *Handler) PasswordRecoveryIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	return h.renderTemplate(rw, ctx)
}

// PasswordRecoverySuccess ...
func (h *Handler) PasswordRecoverySuccess(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
		"message": "web.password_recovery.success_please_check_email",
	})
}

// PasswordRecoveryProcess ...
func (h *Handler) PasswordRecoveryProcess(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	r.ParseForm()

	validationErrorBuilder := lib.NewValidationErrorBuilder()

	passwordRecoveryRequest := request.NewPasswordRecoveryRequestFromForm(r.Form)
	passwordRecoveryRequest.Validate(validationErrorBuilder)

	logger := h.Container.Logger.WithFields(logrus.Fields{
		"email": passwordRecoveryRequest.Email,
	})

	if validationErrorBuilder.HasErrors() {
		rw.WriteHeader(http.StatusBadRequest)
		return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
			"validationErrors": validationErrorBuilder.Errors(),
			"request":          passwordRecoveryRequest,
		})
	}

	err := h.Container.PasswordRecoverer.Start(passwordRecoveryRequest.Email)
	if err != nil {
		switch err {
		case lib.ErrUserNotFound, sql.ErrNoRows:
			logger.Debug("User does not exists. Password cannot be recovered, user will get fake response")
		case lib.ErrPasswordRecovererUserIsNotActive:
			logger.Debug("User is not active. Password cannot be recovered, user will get fake response")
		case lib.ErrPasswordRecovererUserIsNotConfirmed:
			logger.Debug("User is not confirmed. Password cannot be recovered, user will get fake response")
		default:
			return h.renderTemplate500(rw, ctx, err)
		}
	}

	h.redirect(rw, r, "password_recovery_success", http.StatusFound)

	return ctx
}

// PasswordRecoveryConfirmationIndex ...
func (h *Handler) PasswordRecoveryConfirmationIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	var ok bool
	var confirmationToken string
	var userID string

	if confirmationToken, ok = routing.ParamFromContext(ctx, "confirmationToken"); !ok {
		h.Container.Logger.Debug("confirmation token param is missing")
		return h.renderTemplate400(rw, ctx)
	}

	if userID, ok = routing.ParamFromContext(ctx, "userId"); !ok {
		h.Container.Logger.Debug("user id param is missing")
		return h.renderTemplate400(rw, ctx)
	}

	return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
		"confirmation_token": confirmationToken,
		"user_id":            userID,
	})
}

// PasswordRecoveryConfirmationProcess ...
func (h *Handler) PasswordRecoveryConfirmationProcess(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	r.ParseForm()

	var ok bool
	var confirmationToken string
	var userIDParam string

	if confirmationToken, ok = routing.ParamFromContext(ctx, "confirmationToken"); !ok {
		h.Container.Logger.Debug("Confirmation token param is missing.")
		return h.renderTemplate400(rw, ctx)
	}

	if userIDParam, ok = routing.ParamFromContext(ctx, "userId"); !ok {
		h.Container.Logger.Debug("User id param is missing.")
		return h.renderTemplate400(rw, ctx)
	}

	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		h.Container.Logger.Debug("User ID param wrong type.")
		return h.renderTemplate400(rw, ctx)
	}

	logger := h.Container.Logger.WithFields(logrus.Fields{
		"user_id":            userID,
		"confirmation_token": confirmationToken,
	})
	validationErrorBuilder := lib.NewValidationErrorBuilder()

	passwordRecoveryConfirmationRequest := request.NewPasswordRecoveryConfirmationRequestFromForm(r.Form)
	passwordRecoveryConfirmationRequest.Validate(validationErrorBuilder)

	if validationErrorBuilder.HasErrors() {
		rw.WriteHeader(http.StatusBadRequest)
		return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
			"validationErrors":   validationErrorBuilder.Errors(),
			"request":            passwordRecoveryConfirmationRequest,
			"user_id":            userID,
			"confirmation_token": confirmationToken,
		})
	}

	err = h.Container.PasswordRecoverer.Finalize(
		userID,
		confirmationToken,
		passwordRecoveryConfirmationRequest.Password,
	)
	if err != nil {
		switch err {
		case lib.ErrUserNotFound, sql.ErrNoRows:
			logger.Debug("User or PasswordRecovery does not exists. Password cannot be changed, user will get fake response.")
		case lib.ErrPasswordRecovererUserIsNotActive:
			logger.Debug("User is not active. Password cannot be recovered, user will get fake response")
		case lib.ErrPasswordRecovererUserIsNotConfirmed:
			logger.Debug("User is not confirmed. Password cannot be recovered, user will get fake response")
		default:
			return h.renderTemplate500(rw, ctx, err)
		}
	}

	h.redirect(rw, r, "password_recovery_confirmation_success", http.StatusFound)

	return ctx
}

// PasswordRecoveryConfirmationSuccess ...
func (h *Handler) PasswordRecoveryConfirmationSuccess(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
		"message": "web.password_recovery.password_changed",
	})
}
