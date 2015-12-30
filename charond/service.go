package main

import (
	"database/sql"
	stdlog "log"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
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

func initPostgres(connectionString string, retry int, logger log.Logger) *sql.DB {
	var err error
	var attempts int
	var postgres *sql.DB

	// Because of recursion it needs to be checked to not spawn more than one.
	if postgres == nil {
		postgres, err = sql.Open("postgres", connectionString)
		if err != nil {
			sklog.Fatal(logger, err)
		}
	}

	// At this moment connection is not yet established.
	// Ping is required.
	if err := postgres.Ping(); err != nil {
		if attempts > retry {
			sklog.Fatal(logger, err)
		}

		attempts++
		sklog.Error(logger, err)
		time.Sleep(2 * time.Second)

		initPostgres(connectionString, retry, logger)
	} else {
		err = setupDatabase(postgres)
		if err != nil {
			sklog.Fatal(logger, err)
		}
		sklog.Info(logger, "postgres connection has been established", "address", connectionString)
	}

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
