package charon

import (
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func initPostgres(connectionString string, env string, logger log.Logger) (*sql.DB, error) {
	postgres, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if env == EnvironmentTest {
		if err = teardownDatabase(postgres); err != nil {
			return nil, err
		}
		sklog.Info(logger, "database has been cleared upfront")
	}
	err = setupDatabase(postgres)
	if err != nil {
		return nil, err
	}
	sklog.Info(logger, "postgres connection has been established", "address", connectionString)

	return postgres, nil
}

func initMnemosyne(address string, logger log.Logger) (*grpc.ClientConn, mnemosyne.Mnemosyne) {
	if address == "" {
		sklog.Fatal(logger, errors.New("charond: missing mnemosyne address"))
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUserAgent("charon"))
	if err != nil {
		sklog.Fatal(logger, err, "address", address)
	}

	sklog.Info(logger, "rpc connection to mnemosyne has been established", "address", address)

	return conn, mnemosyne.New(conn, mnemosyne.MnemosyneOpts{})
}

func initPasswordHasher(cost int, logger log.Logger) PasswordHasher {
	bh, err := NewBCryptPasswordHasher(cost)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	return bh
}

func initPermissionRegistry(r permissionProvider, permissions Permissions, logger log.Logger) (pr PermissionRegistry) {
	pr = newPermissionRegistry(r)
	created, untouched, removed, err := pr.Register(permissions)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	sklog.Info(logger, "charon permissions has been registered", "created", created, "untouched", untouched, "removed", removed)

	return
}

const (
	MonitoringEnginePrometheus = "prometheus"
)

func initPrometheus(namespace, subsystem string, constLabels stdprometheus.Labels) *monitoring {
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
		enabled: true,
		rpc: monitoringRPC{
			enabled:  true,
			requests: rpcRequests,
			errors:   rpcErrors,
		},
		postgres: monitoringPostgres{
			enabled: true,
			queries: postgresQueries,
			errors:  postgresErrors,
		},
	}
}