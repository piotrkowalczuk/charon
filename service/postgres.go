package service

import (
	"database/sql"
	"errors"

	// ...
	_ "github.com/lib/pq"
)

var (
	// ErrPostgresConnectionFailed ...
	ErrPostgresConnectionFailed = errors.New("service: postgres connection failed")
)

// InitPostgres ...
func InitPostgres(config DBConfig) {
	postgresPool, err := sql.Open("postgres", config.ConnectionString)

	if err != nil {
		Logger.Fatal(ErrPostgresConnectionFailed)
	}

	Logger.Info("Connection do PostgreSQL established successfully.")

	DBPool = postgresPool
}
