package web

import (
	"net/http"

	"golang.org/x/net/context"
)

// MiddlewareFunc ...
type MiddlewareFunc func(*Handler, context.Context, http.ResponseWriter, *http.Request)

// Middlewares ...
type Middlewares []MiddlewareFunc

// NewMiddlewares ...
func NewMiddlewares(fncs ...func(*Handler, context.Context, http.ResponseWriter, *http.Request)) Middlewares {
	middlewares := Middlewares{}

	for _, fn := range fncs {
		middlewares.Add(MiddlewareFunc(fn))
	}

	return middlewares
}

// Add ...
func (m *Middlewares) Add(middleware MiddlewareFunc) {
	*m = append(*m, middleware)
}
