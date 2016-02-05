package main

import (
	"database/sql"
	stdlog "log"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

const (
	loggerAdapterStdOut = "stdout"
	loggerFormatJSON    = "json"
	loggerFormatHumane  = "humane"
	loggerFormatLogFmt  = "logfmt"
)

func initLogger(adapter, format string, level int, context ...interface{}) log.Logger {
	var l log.Logger

	if adapter != loggerAdapterStdOut {
		stdlog.Fatal("service: unsupported logger adapter")
	}

	switch format {
	case loggerFormatHumane:
		l = sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter)
	case loggerFormatJSON:
		l = log.NewJSONLogger(os.Stdout)
	case loggerFormatLogFmt:
		l = log.NewLogfmtLogger(os.Stdout)
	default:
		stdlog.Fatal("charond: unsupported logger format")
	}

	l = log.NewContext(l).With(context...)

	sklog.Info(l, "logger has been initialized successfully", "adapter", adapter, "format", format, "level", level)

	return l
}

func initPostgres(connectionString string, logger log.Logger) *sql.DB {
	postgres, err := sql.Open("postgres", connectionString)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	err = setupDatabase(postgres)
	if err != nil {
		sklog.Fatal(logger, err)
	}
	sklog.Info(logger, "postgres connection has been established", "address", connectionString)

	return postgres
}

func initMnemosyne(address string, logger log.Logger) (*grpc.ClientConn, mnemosyne.Mnemosyne) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		sklog.Fatal(logger, err, "address", address)
	}

	sklog.Info(logger, "rpc connection to mnemosyne has been established", "address", address)

	return conn, mnemosyne.New(conn, mnemosyne.MnemosyneOpts{})
}

func initPasswordHasher(cost int, logger log.Logger) charon.PasswordHasher {
	bh, err := charon.NewBcryptPasswordHasher(cost, logger)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	return bh
}

func initPermissionRegistry(r PermissionRepository, permissions charon.Permissions, logger log.Logger) (pr PermissionRegistry) {
	pr = newPermissionRegistry(r)
	created, untouched, removed, err := pr.Register(permissions)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	sklog.Info(logger, "charon permissions has been registered", "created", created, "untouched", untouched, "removed", removed)

	return
}

const (
	monitoringEnginePrometheus = "prometheus"
)

func initMonitoring(fn func() (*monitoring, error), logger log.Logger) *monitoring {
	m, err := fn()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	return m
}

func initPrometheus(namespace, subsystem string, constLabels stdprometheus.Labels) func() (*monitoring, error) {
	return func() (*monitoring, error) {
		rpcRequests := prometheus.NewCounter(
			stdprometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "rpc_requests_total",
				Help:        "Total number of RPC requests made.",
				ConstLabels: constLabels,
			},
			monitoringRPCLabels,
		)
		rpcErrors := prometheus.NewCounter(
			stdprometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "rpc_errors_total",
				Help:        "Total number of errors that happen during RPC calles.",
				ConstLabels: constLabels,
			},
			monitoringRPCLabels,
		)

		postgresQueries := prometheus.NewCounter(
			stdprometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "postgres_queries_total",
				Help:        "Total number of SQL queries made.",
				ConstLabels: constLabels,
			},
			monitoringPostgresLabels,
		)
		postgresErrors := prometheus.NewCounter(
			stdprometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "postgres_errors_total",
				Help:        "Total number of errors that happen during SQL queries.",
				ConstLabels: constLabels,
			},
			monitoringPostgresLabels,
		)

		return &monitoring{
			rpc: monitoringRPC{
				requests: rpcRequests,
				errors:   rpcErrors,
			},
			postgres: monitoringPostgres{
				queries: postgresQueries,
				errors:  postgresErrors,
			},
		}, nil
	}
}
