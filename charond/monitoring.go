package main

import "github.com/go-kit/kit/metrics"

var (
	monitoringRPCLabels = []string{
		"method",
	}
	monitoringPostgresLabels = []string{
		"query",
	}
)

type monitoring struct {
	rpc      monitoringRPC
	postgres monitoringPostgres
}

type monitoringRPC struct {
	requests metrics.Counter
	errors   metrics.Counter
}

func (mr monitoringRPC) with(f metrics.Field) monitoringRPC {
	return monitoringRPC{
		errors:   mr.errors.With(f),
		requests: mr.requests.With(f),
	}
}

type monitoringPostgres struct {
	queries metrics.Counter
	errors  metrics.Counter
}
