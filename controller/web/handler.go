package web

import (
	"bytes"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// Handler ...
type Handler struct {
	Logger      *logrus.Logger
	DB          *sql.DB
	Middlewares []MiddlewareFunc
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
