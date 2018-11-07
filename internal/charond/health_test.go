package charond

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {
	l := zap.L()
	s := postgresSuite{
		logger: l,
	}
	s.setup(t)

	h := healthHandler{
		postgres: s.db,
		logger:   l,
	}

	rw := httptest.NewRecorder()
	r := &http.Request{}
	h.ServeHTTP(rw, r)
	if rw.Code != http.StatusOK {
		t.Errorf("wrong status code, expected %d but got %d", http.StatusOK, rw.Code)
	}

	s.teardown(t)

	rw = httptest.NewRecorder()
	h.ServeHTTP(rw, r)
	if rw.Code != http.StatusServiceUnavailable {
		t.Errorf("wrong status code, expected %d but got %d", http.StatusServiceUnavailable, rw.Code)
	}
}
