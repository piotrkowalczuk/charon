package web

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-soa/auth/lib/security"
	"github.com/go-soa/auth/repository"
	"github.com/go-soa/auth/service"
	"golang.org/x/net/context"
)

// Handler ...
type Handler struct {
	Logger         *logrus.Logger
	DB             *sql.DB
	RM             repository.Manager
	PasswordHasher security.PasswordHasher
	Middlewares    []MiddlewareFunc
	TmplName       string
	Tmpl           *template.Template
	Mailer         service.Mailer
}

// ServeHTTP ...
func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	startedAt := time.Now()
	wrw := &ResponseWriter{
		ResponseWriter: rw,
	}
	ctx := context.TODO()

	for _, middleware := range h.Middlewares {
		middleware(h, ctx, wrw, r)
	}

	if wrw.StatusCode == 0 {
		wrw.StatusCode = http.StatusOK
	}

	h.logRequest(wrw, r, startedAt)
}

func (h *Handler) logRequest(wrw *ResponseWriter, r *http.Request, startedAt time.Time) {
	b := bytes.NewBufferString("[")
	b.WriteString(r.Method)
	b.WriteString("] ")
	b.WriteString("[")
	b.WriteString(strconv.FormatInt(int64(wrw.StatusCode), 10))
	b.WriteString("] ")
	b.WriteString(r.URL.Path)
	b.WriteString(" ")
	b.WriteString(time.Since(startedAt).String())

	h.Logger.Info(b.String())
}
