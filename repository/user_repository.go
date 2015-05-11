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
	// ErrUserNotFound  ...
	ErrUserNotFound = errors.New("repository: user not found")
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

func (ur *UserRepository) RegistrationConfirmation(userID int64, confirmationToken string) error {
	query := `
		UPDATE auth_user
		SET is_confirmed = true, updated_at = NOW()
		WHERE is_confirmed = false AND id = $1 AND confirmation_token = $2;
	`

	result, err := ur.db.Exec(
		query,
		userID,
		confirmationToken,
	)

	if err != nil {
		return mapKnownErrors(userKnownErrors, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return mapKnownErrors(userKnownErrors, err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
