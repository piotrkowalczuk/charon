package charon

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestInitPrometheus(t *testing.T) {
	monitoring := initPrometheus("namespace", "subsystem", prometheus.Labels{"server": "travis-ci"})
	if monitoring == nil {
		t.Fatalf("nil monitoring")
	}
}
