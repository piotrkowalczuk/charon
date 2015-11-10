package service

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon/charond/lib"
)

// Singleton instance of lib.RepositoryManager.
var RepositoryManager lib.RepositoryManager

// InitRepositoryManager ...
func InitRepositoryManager(db *sql.DB) {
	rm := lib.RepositoryManager{
		User:             lib.NewUserRepository(db),
		PasswordRecovery: lib.NewPasswordRecoveryRepository(db),
	}

	RepositoryManager = rm
}
