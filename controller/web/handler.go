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
	"golang.org/x/net/context"
	"github.com/go-soa/auth/mail"
)

// ServiceContainer ...
type ServiceContainer struct {
	Logger         *logrus.Logger
	DB             *sql.DB
	RM             repository.Manager
	PasswordHasher security.PasswordHasher
	Mailer         *mail.Mail
	Templates      *template.Template
}

// Handler ...
type Handler struct {
	TemplateName string
	Middlewares  []MiddlewareFunc
	Container    ServiceContainer
}

// NewHandler ...
func NewHandler(templateName string, middlewares []MiddlewareFunc, container ServiceContainer) *Handler {
	return &Handler{
		TemplateName: templateName,
		Middlewares:  middlewares,
		Container:    container,
	}
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

	h.Container.Logger.Info(b.String())
}
