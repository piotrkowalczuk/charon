package repository

import (
	"database/sql"
	"errors"

	"github.com/go-soa/charon/model"
)

const (
	userUniqueConstraintViolationUsernameErrorMessage = "pq: duplicate key value violates unique constraint \"auth_user_username_key\""
)

var (
	// ErrUserUniqueConstraintViolationUsername ...
	ErrUserUniqueConstraintViolationUsername = errors.New(userUniqueConstraintViolationUsernameErrorMessage)
	userKnownErrors                          = map[string]error{
		userUniqueConstraintViolationUsernameErrorMessage: ErrUserUniqueConstraintViolationUsername,
	}
)

// UserRepository ...
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository ...
func NewUserRepository(dbPool *sql.DB) (repository *UserRepository) {
	repository = &UserRepository{dbPool}

	return
}

// Insert ...
func (ur *UserRepository) Insert(user *model.User) (sql.Result, error) {
	query := `
		INSERT INTO auth_user (
			password, username, first_name, last_name, is_active, is_staff,
			is_superuser, is_confirmed, confirmation_token, last_login_at,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	result, err := ur.db.Exec(
		query,
		user.Password,
		user.Username,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.IsStaff,
		user.IsSuperuser,
		user.IsConfirmed,
		user.ConfirmationToken,
		user.LastLoginAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return result, mapKnownErrors(userKnownErrors, err)
}
