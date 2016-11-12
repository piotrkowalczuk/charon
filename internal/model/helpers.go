package model

import (
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
)

func untouched(given, created, removed int64) int64 {
	switch {
	case given < 0:
		return -1
	case given == 0:
		return -2
	case given < created:
		return 0
	default:
		return given - created
	}
}

func initPostgres(address string, test bool, logger log.Logger) (*sql.DB, error) {
	postgres, err := sql.Open("postgres", address)
	if err != nil {
		return nil, err
	}

	if test {
		if err = teardownDatabase(postgres); err != nil {
			return nil, err
		}
		sklog.Info(logger, "database has been cleared upfront")
	}
	err = setupDatabase(postgres)
	if err != nil {
		return nil, err
	}

	sklog.Info(logger, "postgres connection has been established", "address", address)

	return postgres, nil
}
