package service

import (
	"database/sql"
	"time"

	// ...
	_ "github.com/lib/pq"
)

// InitPostgres ...
func InitPostgres(config DBConfig) {
	var err error

	// Because of recursion it needs to be checked to not spawn more than one.
	if DBPool == nil {
		DBPool, err = sql.Open("postgres", config.ConnectionString)
		if err != nil {
			Logger.Fatal(err)
		}
	}

	// At this moment connection is not yet established.
	// Ping is required.
	if err := DBPool.Ping(); err != nil {
		Logger.Error(err)
		time.Sleep(2 * time.Second)

		InitPostgres(config)
	} else {
		Logger.Info("Connection do PostgreSQL established successfully.")
	}
}
