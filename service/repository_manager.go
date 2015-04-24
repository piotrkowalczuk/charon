package service

import (
	"database/sql"

	"github.com/go-soa/auth/repository"
)

// Singleton instance of repository.Manager.
var RepositoryManager repository.Manager

// InitRepositoryManager ...
func InitRepositoryManager(db *sql.DB) {
	rm := repository.Manager{
		User: repository.NewUserRepository(db),
	}

	RepositoryManager = rm
}
