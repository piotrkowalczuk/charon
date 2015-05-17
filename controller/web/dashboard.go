package web

import (
	"net/http"

	mnemosyne "github.com/go-soa/mnemosyne/lib"
	"golang.org/x/net/context"
)

// RegistrationIndex ...
func (h *Handler) DashboardIndex(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	session, ok := ctx.Value("session").(mnemosyne.Session)
	if !ok {
		return h.renderTemplate(rw, ctx)
	}

	return h.renderTemplateWithData(rw, ctx, map[string]interface{}{
		"session": session.Data,
	})
}
