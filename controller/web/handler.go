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
	"github.com/lib/pq"
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
func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	startedAt := time.Now()
	ctx := routing.NewParamsContext(context.Background(), ps)
	wrw := &ResponseWriter{
		ResponseWriter: rw,
	}

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

	router.Handle(method, pattern, httprouter.Handle(h.ServeHTTP))

	h.Container.Logger.WithFields(logrus.Fields{
		"name":    routeName,
		"method":  method,
		"pattern": pattern,
	}).Info("Route has been registered successfully.")
}

func (h *Handler) sendErrorWithStatus(rw http.ResponseWriter, err error, status int) {
	switch e := err.(type) {
	case *pq.Error:
		h.Container.Logger.WithFields(logrus.Fields{
			"severity":   e.Severity,
			"code":       e.Code,
			"detail":     e.Detail,
			"hint":       e.Hint,
			"position":   e.Position,
			"table":      e.Table,
			"constraint": e.Constraint,
		}).Error(e.Message)
	default:
		h.Container.Logger.Error(e)
	}

	http.Error(rw, err.Error(), status)
}

func (h *Handler) sendError500(rw http.ResponseWriter, err error) {
	h.sendErrorWithStatus(rw, err, http.StatusInternalServerError)
}

func (h *Handler) sendError400(rw http.ResponseWriter, err error) {
	h.sendErrorWithStatus(rw, err, http.StatusBadRequest)
}

func (h *Handler) renderTemplate(rw http.ResponseWriter) {
	h.renderTemplateWithData(rw, nil)
}

func (h *Handler) renderTemplateWithData(rw http.ResponseWriter, data interface{}) {
	err := h.Container.Templates.ExecuteTemplate(rw, h.Name, data)
	if err != nil {
		h.sendError500(rw, err)
		return
	}
}

func (h *Handler) renderTemplateWithStatus(rw http.ResponseWriter, status int) {
	var templateName string

	switch {
	case status >= 400 && status < 500:
		templateName = "400"
	case status == 404:
		templateName = "404"
	default:
		templateName = "500"
	}

	err := h.Container.Templates.ExecuteTemplate(rw, templateName, map[string]string{
		"status": strconv.FormatInt(int64(status), 10),
	})
	if err != nil {
		h.sendError500(rw, err)
		return
	}
}

func (h *Handler) renderTemplate404(rw http.ResponseWriter) {
	h.renderTemplateWithStatus(rw, http.StatusNotFound)
}

func (h *Handler) renderTemplate400(rw http.ResponseWriter) {
	h.renderTemplateWithStatus(rw, http.StatusBadRequest)
}

func (h *Handler) renderTemplate500(rw http.ResponseWriter) {
	h.renderTemplateWithStatus(rw, http.StatusInternalServerError)
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
