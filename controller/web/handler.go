package web

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-soa/charon/lib/routing"
	"github.com/go-soa/charon/lib/security"
	"github.com/go-soa/charon/mail"
	"github.com/go-soa/charon/repository"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// ServiceContainer ...
type ServiceContainer struct {
	Logger             *logrus.Logger
	DB                 *sql.DB
	ConfirmationMailer mail.Sender
	Templates          *template.Template
	RM                 repository.Manager
	PasswordHasher     security.PasswordHasher
	Routes             routing.Routes
	URLGenerator       routing.URLGenerator
}

type HandlerOpts struct {
	Name        string
	Method      string
	Middlewares []MiddlewareFunc
	Container   ServiceContainer
}

// Handler ...
type Handler struct {
	Name        string
	Method      string
	middlewares []MiddlewareFunc
	Container   ServiceContainer
}

// NewHandler ...
func NewHandler(options HandlerOpts) *Handler {
	return &Handler{
		Name:        options.Name,
		Method:      options.Method,
		middlewares: options.Middlewares,
		Container:   options.Container,
	}
}

// RouteName ...
func (h *Handler) RouteName() routing.RouteName {
	return routing.RouteName(h.Name)
}

// ServeHTTP ...
func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	startedAt := time.Now()
	wrw := &ResponseWriter{
		ResponseWriter: rw,
	}
	ctx := context.TODO()

	for _, middleware := range h.middlewares {
		middleware(h, ctx, wrw, r)
	}

	if wrw.StatusCode == 0 {
		wrw.StatusCode = http.StatusOK
	}

	h.logRequest(wrw, r, startedAt)
}

// Register ...
func (h *Handler) Register(router *httprouter.Router) {
	method := h.Method
	routeName := h.RouteName()
	pattern := h.Container.Routes.GetPattern(routeName).String()

	router.Handler(method, pattern, h)

	h.Container.Logger.WithFields(logrus.Fields{
		"name":    routeName,
		"method":  method,
		"pattern": pattern,
	}).Info("Route has been registered successfully.")
}

func (h *Handler) sendErrorWithStatus(rw http.ResponseWriter, err error, status int) {
	h.Container.Logger.Error(err)
	http.Error(rw, err.Error(), status)
}

func (h *Handler) sendError500(rw http.ResponseWriter, err error) {
	h.sendErrorWithStatus(rw, err, http.StatusInternalServerError)
}

func (h *Handler) renderTemplate(rw http.ResponseWriter) {
	h.renderTemplateWithData(rw, nil)
}

func (h *Handler) renderTemplateWithData(rw http.ResponseWriter, data interface{}) {
	err := h.Container.Templates.ExecuteTemplate(rw, h.Name, data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
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
