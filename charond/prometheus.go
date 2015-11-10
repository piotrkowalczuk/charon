package main

import (
	pmetrics "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	monitoringEnginePrometheus = "prometheus"
)

func initPrometheus(namespace, subsystem string, constLabels prometheus.Labels) func() (*monitoring, error) {
	return func() (*monitoring, error) {
		rpcRequests := pmetrics.NewCounter(
			prometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "rpc_requests_total",
				Help:        "Total number of RPC requests made.",
				ConstLabels: constLabels,
			},
			monitoringRPCLabels,
		)
		rpcErrors := pmetrics.NewCounter(
			prometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "rpc_errors_total",
				Help:        "Total number of errors that happen during RPC calles.",
				ConstLabels: constLabels,
			},
			monitoringRPCLabels,
		)

		postgresQueries := pmetrics.NewCounter(
			prometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "postgres_queries_total",
				Help:        "Total number of SQL queries made.",
				ConstLabels: constLabels,
			},
			monitoringPostgresLabels,
		)
		postgresErrors := pmetrics.NewCounter(
			prometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "postgres_errors_total",
				Help:        "Total number of errors that happen during SQL queries.",
				ConstLabels: constLabels,
			},
			monitoringPostgresLabels,
		)

		m := &monitoring{}
		m.rpc.requests = rpcRequests
		m.rpc.errors = rpcErrors
		m.postgres.queries = postgresQueries
		m.postgres.errors = postgresErrors

		return m, nil
	}
}
