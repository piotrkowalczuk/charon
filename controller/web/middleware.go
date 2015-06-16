package web

import (
	"net/http"

	"golang.org/x/net/context"
)

// MiddlewareFunc ...
type MiddlewareFunc func(*Handler, context.Context, http.ResponseWriter, *http.Request) context.Context

// Middlewares ...
type Middlewares []MiddlewareFunc

// NewMiddlewares ...
func NewMiddlewares(fns ...func(*Handler, context.Context, http.ResponseWriter, *http.Request) context.Context) Middlewares {
	middlewares := Middlewares{}

	for _, fn := range fns {
		middlewares.Add(MiddlewareFunc(fn))
	}

	return middlewares
}

// Add ...
func (m *Middlewares) Add(middleware MiddlewareFunc) {
	*m = append(*m, middleware)
}
