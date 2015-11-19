package main

import (
	"database/sql"
	"time"
)

const (
	// UserConfirmationTokenUsed is a value that is used when confirmation token was already used.
	UserConfirmationTokenUsed = "!"
	sqlSchemaUser             = `
		CREATE TABLE IF NOT EXISTS (
			id serial PRIMARY KEY,
			password character varying(128) NOT NULL,
			username character varying(75) NOT NULL,
			first_name character varying(45) NOT NULL,
			last_name character varying(45) NOT NULL,
			is_superuser boolean NOT NULL,
			is_active boolean NOT NULL,
			is_staff boolean NOT NULL,
			is_confirmed boolean NOT NULL,
			confirmation_token character varying(122) NOT NULL,
			last_login_at timestamp with time zone NOT NULL,
			created_at timestamp with time zone NOT NULL,
			updated_at timestamp with time zone NOT NULL
		)
	`
)

type userEntity struct {
	ID                int64
	Password          string
	Username          string
	FirstName         string
	LastName          string
	IsActive          bool
	IsStaff           bool
	IsSuperuser       bool
	IsConfirmed       bool
	ConfirmationToken string
	LastLoginAt       *time.Time
	CreatedAt         *time.Time
	CreatedBy         sql.NullInt64
	UpdatedAt         *time.Time
	UpdatedBy         sql.NullInt64
}

// String return concatenated first and last name of the user.
func (u *userEntity) String() string {
	return u.FirstName + " " + u.LastName
}

// UserRepository ...
type UserRepository interface {
	Create(username, password, firstName, lastName, confirmationToken string, isSuperuser, isStaff, isActive, isConfirmed bool) (*userEntity, error)
	CreateSuperuser(username, password, firstName, lastName string) (*userEntity, error)
	// Count retrieves number of all users.
	Count() (int64, error)
	UpdateLastLoginAt(id int64) error
	ChangePassword(id int64, password string) error
	FindOneByID(id int64) (*userEntity, error)
	FindOneByUsername(username string) (*userEntity, error)
	RegistrationConfirmation(id int64, confirmationToken string) error
}

type userRepository struct {
	db *sql.DB
}

func newUserRepository(dbPool *sql.DB) *userRepository {
	return &userRepository{
		db: dbPool,
	}
}

// Create implements UserRepository interface.
func (ur *userRepository) Create(username, password, firstName, lastName, confirmationToken string, isSuperuser, isStaff, isActive, isConfirmed bool) (*userEntity, error) {
	entity := &userEntity{
		Username:          username,
		Password:          password,
		FirstName:         firstName,
		LastName:          lastName,
		ConfirmationToken: confirmationToken,
		IsSuperuser:       isSuperuser,
		IsStaff:           isStaff,
		IsActive:          isActive,
		IsConfirmed:       isConfirmed,
	}
	err := ur.insert(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// CreateSuperuser implements UserRepository interface.
func (ur *userRepository) CreateSuperuser(username, password, firstName, lastName string) (*userEntity, error) {
	return ur.Create(username, password, firstName, lastName, UserConfirmationTokenUsed, true, false, true, true)
}

// Count implements UserRepository interface.
func (ur *userRepository) Count() (n int64, err error) {
	err = ur.db.QueryRow("SELECT COUNT(*) FROM charon.user").Scan(&n)

	return
}

func (ur *userRepository) insert(e *userEntity) error {
	query := `
		INSERT INTO charon.user (
			password, username, first_name, last_name, is_active, is_staff,
			is_superuser, is_confirmed, confirmation_token,
			created_at, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), $10)
		RETURNING id, created_at
	`
	return ur.db.QueryRow(
		query,
		e.Password,
		e.Username,
		e.FirstName,
		e.LastName,
		e.IsActive,
		e.IsStaff,
		e.IsSuperuser,
		e.IsConfirmed,
		e.ConfirmationToken,
		e.CreatedBy,
	).Scan(&e.ID, &e.CreatedAt)
}

func (ur *userRepository) RegistrationConfirmation(userID int64, confirmationToken string) error {
	query := `
		UPDATE charon.user
		SET is_confirmed = true, is_active = true, updated_at = NOW(), confirmation_token = $1
		WHERE is_confirmed = false AND id = $2 AND confirmation_token = $3;
	`

	result, err := ur.db.Exec(
		query,
		UserConfirmationTokenUsed,
		userID,
		confirmationToken,
	)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}

// ChangePassword ...
func (ur *userRepository) ChangePassword(userID int64, password string) error {
	query := `
		UPDATE charon.user
		SET password = $2, updated_at = NOW()
		WHERE id = $1;
	`

	result, err := ur.db.Exec(
		query,
		userID,
		password,
	)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()

	return err
}

// FindOneByUsername ...
func (ur *userRepository) FindOneByUsername(username string) (*userEntity, error) {
	return ur.findOneBy("username", username)
}

// FindOneByID ...
func (ur *userRepository) FindOneByID(id int64) (*userEntity, error) {
	return ur.findOneBy("id", id)
}

func (ur *userRepository) findOneBy(fieldName string, value interface{}) (*userEntity, error) {
	query := `
		SELECT id, password, username, first_name, last_name, is_active, is_staff,
			is_superuser, is_confirmed, confirmation_token, last_login_at,
			created_at, updated_at
		FROM charon.user
		WHERE ` + fieldName + ` = $1
		LIMIT 1
	`

	user := &userEntity{}
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

// UpdateLastLoginAt implements UserRepository interface.
func (ur *userRepository) UpdateLastLoginAt(userID int64) error {
	query := `
		UPDATE charon.user
		SET last_login_at = NOW()
		WHERE id = $1;
	`

	result, err := ur.db.Exec(
		query,
		userID,
	)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()

	return err
}
