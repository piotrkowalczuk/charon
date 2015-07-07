package lib

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
	ErrUserNotFound = errors.New("charon: user not found")
	// ErrUserUniqueConstraintViolationUsername ...
	ErrUserUniqueConstraintViolationUsername = errors.New(userUniqueConstraintViolationUsernameErrorMessage)
	userKnownErrors                          = map[string]error{
		userUniqueConstraintViolationUsernameErrorMessage: ErrUserUniqueConstraintViolationUsername,
	}
)

// UserRepository ...
type UserRepository interface {
	Insert(*model.User) error
	UpdateLastLoginAt(int64) error
	ChangePassword(int64, string) error
	FindOneByID(int64) (*model.User, error)
	FindOneByUsername(string) (*model.User, error)
	RegistrationConfirmation(int64, string) error
}

// userRepository ...
type userRepository struct {
	db *sql.DB
}

// NewUserRepository ...
func NewUserRepository(dbPool *sql.DB) (repository *userRepository) {
	repository = &userRepository{dbPool}

	return
}

// Insert ...
func (ur *userRepository) Insert(user *model.User) error {
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

func (ur *userRepository) RegistrationConfirmation(userID int64, confirmationToken string) error {
	query := `
		UPDATE charon_user
		SET is_confirmed = true, is_active = true, updated_at = NOW(), confirmation_token = $1
		WHERE is_confirmed = false AND id = $2 AND confirmation_token = $3;
	`

	result, err := ur.db.Exec(
		query,
		model.UserConfirmationTokenUsed,
		userID,
		confirmationToken,
	)
	if err != nil {
		return mapKnownErrors(userKnownErrors, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// ChangePassword ...
func (ur *userRepository) ChangePassword(userID int64, password string) error {
	query := `
		UPDATE charon_user
		SET password = $2, updated_at = NOW()
		WHERE id = $1;
	`

	result, err := ur.db.Exec(
		query,
		userID,
		password,
	)
	if err != nil {
		return mapKnownErrors(userKnownErrors, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// FindOneByUsername ...
func (ur *userRepository) FindOneByUsername(username string) (*model.User, error) {
	return ur.FindOneBy("username", username)
}

// FindOneByID ...
func (ur *userRepository) FindOneByID(id int64) (*model.User, error) {
	return ur.FindOneBy("id", id)
}

func (ur *userRepository) FindOneBy(fieldName string, value interface{}) (*model.User, error) {
	query := `
		SELECT id, password, username, first_name, last_name, is_active, is_staff,
			is_superuser, is_confirmed, confirmation_token, last_login_at,
			created_at, updated_at
		FROM charon_user
		WHERE ` + fieldName + ` = $1
		LIMIT 1
	`

	user := &model.User{}
	err := ur.db.QueryRow(query, value).Scan(
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

func (ur *userRepository) UpdateLastLoginAt(userID int64) error {
	query := `
		UPDATE charon_user
		SET last_login_at = NOW()
		WHERE id = $1;
	`

	result, err := ur.db.Exec(
		query,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return mapKnownErrors(userKnownErrors, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
