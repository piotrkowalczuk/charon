package repository

import (
	"database/sql"
	"errors"

	"github.com/go-soa/charon/model"
)

const (
	userUniqueConstraintViolationUsernameErrorMessage = "pq: duplicate key value violates unique constraint \"charon_user_username_key\""
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

// UserProvider ...
type UserProvider interface {
	FindByUsername(string) *model.User
}

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
func (ur *UserRepository) Insert(user *model.User) error {
	query := `
		INSERT INTO charon_user (
			password, username, first_name, last_name, is_active, is_staff,
			is_superuser, is_confirmed, confirmation_token, last_login_at,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	err := ur.db.QueryRow(
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
	).Scan(&user.ID)

	return mapKnownErrors(userKnownErrors, err)
}

func (ur *UserRepository) RegistrationConfirmation(userID int64, confirmationToken string) error {
	query := `
		UPDATE charon_user
		SET is_confirmed = true, is_active = true, updated_at = NOW()
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

func (ur *UserRepository) FindByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, password, username, first_name, last_name, is_active, is_staff,
			is_superuser, is_confirmed, confirmation_token, last_login_at,
			created_at, updated_at
		FROM charon_user
		WHERE username = $1
		LIMIT 1
	`

	user := &model.User{}
	err := ur.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Password,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.IsStaff,
		&user.IsSuperuser,
		&user.IsConfirmed,
		&user.ConfirmationToken,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
