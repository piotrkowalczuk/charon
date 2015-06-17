package web

import (
	"bytes"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"net/rpc"

	"github.com/Sirupsen/logrus"
	"github.com/go-soa/charon/lib/routing"
	"github.com/go-soa/charon/lib/security"
	"github.com/go-soa/charon/mail"
	"github.com/go-soa/charon/repository"
	"github.com/go-soa/charon/service"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"golang.org/x/net/context"
)

// ServiceContainer ...
type ServiceContainer struct {
	Config             service.AppConfig
	Logger             *logrus.Logger
	DB                 *sql.DB
	ConfirmationMailer mail.Sender
	RM                 repository.Manager
	PasswordHasher     security.PasswordHasher
	Routes             routing.Routes
	URLGenerator       routing.URLGenerator
	TemplateManager    *service.TemplateManager
	Mnemosyne          *rpc.Client
}

// HandlerOpts ...
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
	cancel      context.CancelFunc
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
	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel
	ctx = routing.NewParamsContext(ctx, ps)
	wrw := &ResponseWriter{
		ResponseWriter: rw,
	}

MiddlewaresLoop:
	for _, middleware := range h.middlewares {
		ctx = middleware(h, ctx, wrw, r)

		select {
		case <-ctx.Done():
			break MiddlewaresLoop
		default:
			continue
		}
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
	}).Info("Web view has been registered successfully.")
}

func (h *Handler) logError(err error) {
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
}

func (h *Handler) renderTemplate(rw http.ResponseWriter, ctx context.Context) context.Context {
	return h.renderTemplateWithData(rw, ctx, nil)
}

func (h *Handler) renderTemplateWithData(rw http.ResponseWriter, ctx context.Context, data interface{}) context.Context {
	err := h.Container.TemplateManager.GetForWeb(rw, h.Name, data)

	if err != nil {
		h.cancel()
		h.logError(err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}

	return ctx
}

func (h *Handler) renderTemplateWithStatus(rw http.ResponseWriter, ctx context.Context, status int) context.Context {
	h.cancel()
	var templateName string

	switch {
	case status == 404:
		templateName = "404"
		break
	case status >= 400 && status < 500:
		templateName = "400"
		break
	default:
		templateName = "500"
	}

	err := h.Container.TemplateManager.GetForWeb(rw, templateName, map[string]interface{}{
		"status": strconv.FormatInt(int64(status), 10),
	})
	if err != nil {
		h.logError(err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return ctx
	}

	return ctx
}

func (h *Handler) renderTemplate404(rw http.ResponseWriter, ctx context.Context) context.Context {
	return h.renderTemplateWithStatus(rw, ctx, http.StatusNotFound)
}

func (h *Handler) renderTemplate403(rw http.ResponseWriter, ctx context.Context) context.Context {
	return h.renderTemplateWithStatus(rw, ctx, http.StatusForbidden)
}

func (h *Handler) renderTemplate400(rw http.ResponseWriter, ctx context.Context) context.Context {
	return h.renderTemplateWithStatus(rw, ctx, http.StatusBadRequest)
}

func (h *Handler) renderTemplate500(rw http.ResponseWriter, ctx context.Context, err error) context.Context {
	h.logError(err)
	return h.renderTemplateWithStatus(rw, ctx, http.StatusInternalServerError)
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
