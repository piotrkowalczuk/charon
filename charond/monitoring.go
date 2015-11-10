package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/piotrkowalczuk/sklog"
)

var (
	monitor             *monitoring
	monitoringRPCLabels = []string{
		"method",
	}
	monitoringPostgresLabels = []string{
		"query",
	}
)

type monitoring struct {
	rpc struct {
		requests metrics.Counter
		errors   metrics.Counter
	}
	postgres struct {
		queries metrics.Counter
		errors  metrics.Counter
	}
}

func initMonitoring(fn func() (*monitoring, error), logger log.Logger) {
	m, err := fn()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	monitor = m
}
