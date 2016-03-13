// +build unit,!postgres,!e2e

package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestInitPrometheus(t *testing.T) {
	_, err := initPrometheus("namespace", "subsystem", prometheus.Labels{"server": "travis-ci"})()
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}
