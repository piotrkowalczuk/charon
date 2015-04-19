package web

import "net/http"

// ResponseWriter ...
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader ...
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
