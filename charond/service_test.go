package charond

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestInitPrometheus(t *testing.T) {
	monitoring := initPrometheus("namespace", true, prometheus.Labels{"server": "travis-ci"})
	if monitoring == nil {
		t.Fatalf("nil monitoring")
	}
}
