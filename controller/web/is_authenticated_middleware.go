package web

import (
	"net/http"

	mnemosyne "github.com/go-soa/mnemosyne/lib"
	"golang.org/x/net/context"
)

// IsAuthenticatedMiddleware ...
func (h *Handler) IsAuthenticatedMiddleware(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	cookie, err := r.Cookie("sid")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return h.renderTemplate400(rw, ctx)
		default:
			return h.renderTemplate500(rw, ctx, err)
		}
	}

	session := mnemosyne.Session{}
	err = h.Container.Mnemosyne.Call("Store.Get", mnemosyne.SessionID(cookie.Value), &session)
	if err != nil {
		switch err.Error() {
		case mnemosyne.ErrSessionNotFound.Error():
			return h.renderTemplate403(rw, ctx)
		default:
			return h.renderTemplate500(rw, ctx, err)
		}
	}

	// TODO(piotr): current user status need to be checked (is_active, is_confirmed etc)
	return context.WithValue(ctx, "session", session)
}
