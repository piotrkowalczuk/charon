package main

import (
	"database/sql"
	"time"

	// ...

	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/sklog"
)

var postgres *sql.DB

func initPostgres(connectionString string, retry int, logger log.Logger) {
	var err error
	var attempts int

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
		sklog.Info(logger, "connection do postgres established successfully")
	}
}
