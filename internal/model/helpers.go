package model

import (
	"database/sql"

	"go.uber.org/zap"
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

func initPostgres(address string, test bool, logger *zap.Logger) (*sql.DB, error) {
	postgres, err := sql.Open("postgres", address)
	if err != nil {
		return nil, err
	}

	if test {
		if err = teardownDatabase(postgres); err != nil {
			return nil, err
		}
		logger.Info("database has been cleared upfront")
	}
	err = setupDatabase(postgres)
	if err != nil {
		return nil, err
	}

	logger.Info("postgres connection has been established", zap.String("address", address))

	return postgres, nil
}
