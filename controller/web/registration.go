package web

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

// RegistrationGET ...
func (h *Handler) RegistrationGET(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Welcome!\n")
}
