package service

import (
	"database/sql"
	"errors"
)

const (
	// DBAdapterPostgres ...
	DBAdapterPostgres = "postgres"
	// DBAdapterMySQL ...
	DBAdapterMySQL = "mysql"
)

var (
	// DBPool ...
	DBPool = &sql.DB{}
	// ErrDatabaseAdapterNotSupported ...
	ErrDatabaseAdapterNotSupported = errors.New("service: database adapter not supported")
)

// DBConfig ...
type DBConfig struct {
	Adapter          string `xml:"adapter"`
	ConnectionString string `xml:"connection-string"`
}

// InitDB ...
func InitDB(config DBConfig) {
	switch config.Adapter {
	case DBAdapterPostgres:
		InitPostgres(config)
	case DBAdapterMySQL:
		Logger.Fatal(ErrDatabaseAdapterNotSupported)
	default:
		Logger.Fatal(ErrDatabaseAdapterNotSupported)
	}
}
