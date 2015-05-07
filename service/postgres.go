package service

import (
	"database/sql"
	"time"

	// ...
	_ "github.com/lib/pq"
)

// PostgresPool ...
var PostgresPool *sql.DB

// InitPostgres ...
func InitPostgres(config DBConfig) {
	var err error

	// Because of recursion it needs to be checked to not spawn more than one.
	if PostgresPool == nil {
		PostgresPool, err = sql.Open("postgres", config.ConnectionString)
		if err != nil {
			Logger.Fatal(err)
		}
	}

	// At this moment connection is not yet established.
	// Ping is required.
	if err := PostgresPool.Ping(); err != nil {
		Logger.Error(err)
		time.Sleep(2 * time.Second)

		InitPostgres(config)
	} else {
		Logger.Info("Connection do PostgreSQL established successfully.")
	}
}
