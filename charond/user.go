package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/pqcnstr"
	"github.com/piotrkowalczuk/protot"
)

const (
	// UserConfirmationTokenUsed is a value that is used when confirmation token was already used.
	UserConfirmationTokenUsed                     = "!"
	sqlCnstrPrimaryKeyUser     pqcnstr.Constraint = "charon.user_pkey"
	sqlCnstrUniqueUserUsername pqcnstr.Constraint = "charon.user_username_key"
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
func (ue *userEntity) String() string {
	return ue.FirstName + " " + ue.LastName
}

func (ue *userEntity) subjectID() string {
	return strconv.FormatInt(ue.ID, 10)
}

// Message allocates new corresponding protobuf message.
func (ue *userEntity) Message() *charon.User {
	var (
		createdAt *protot.Timestamp
		updatedAt *protot.Timestamp
	)

	if ue.CreatedAt != nil {
		createdAt = protot.TimeToTimestamp(*ue.CreatedAt)
	}

	if ue.UpdatedAt != nil {
		createdAt = protot.TimeToTimestamp(*ue.UpdatedAt)
	}
	return &charon.User{
		Id:          ue.ID,
		Username:    ue.Username,
		FirstName:   ue.FirstName,
		LastName:    ue.LastName,
		IsSuperuser: ue.IsSuperuser,
		IsActive:    ue.IsActive,
		IsStaff:     ue.IsStaff,
		IsConfirmed: ue.IsConfirmed,
		CreatedAt:   createdAt,
		CreatedBy:   ue.CreatedBy.Int64,
		UpdatedAt:   updatedAt,
		UpdatedBy:   ue.UpdatedBy.Int64,
	}
}

// UserRepository ...
type UserRepository interface {
	Create(username, password, firstName, lastName, confirmationToken string, isSuperuser, isStaff, isActive, isConfirmed bool) (*userEntity, error)
	CreateSuperuser(username, password, firstName, lastName string) (*userEntity, error)
	// Count retrieves number of all users.
	Count() (int64, error)
	UpdateLastLoginAt(id int64) error
	ChangePassword(id int64, password string) error
	Find(offset, limit *protot.NilInt64) ([]*userEntity, error)
	FindOneByID(id int64) (*userEntity, error)
	FindOneByUsername(username string) (*userEntity, error)
	DeleteOneByID(id int64) (int64, error)
	UpdateOneByID(id int64, username, securePassword, firstName, lastName *protot.NilString, isSuperuser, isActive, isStaff, isConfirmed *protot.NilBool) (*userEntity, error)
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

	var user userEntity
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

	return &user, nil
}

// Find implements UserRepository interface.
func (ur *userRepository) Find(offset, limit *protot.NilInt64) ([]*userEntity, error) {
	query := `
		SELECT id, password, username, first_name, last_name, is_active, is_staff,
			is_superuser, is_confirmed, confirmation_token, last_login_at,
			created_at, updated_at
		FROM charon.user
		OFFSET $1
		LIMIT $2
	`

	if offset == nil && !offset.Valid {
		offset = &protot.NilInt64{Int64: 0, Valid: true}
	}

	if limit == nil || !limit.Valid {
		limit = &protot.NilInt64{Int64: 10, Valid: true}
	}

	rows, err := ur.db.Query(query, offset.Int64, limit.Int64)
	if err != nil {
		return nil, err
	}

	users := []*userEntity{}
	for rows.Next() {
		var user userEntity
		err = rows.Scan(
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
		users = append(users, &user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
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

func (ur *userRepository) DeleteOneByID(id int64) (int64, error) {
	query := `
		DELETE FROM charon.user
		WHERE id = $1
	`

	res, err := ur.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (ur *userRepository) UpdateOneByID(id int64, username, securePassword, firstName, lastName *protot.NilString, isSuperuser, isActive, isStaff, isConfirmed *protot.NilBool) (*userEntity, error) {
	keys := make([]string, 0, 8)
	values := make([]interface{}, 0, 9)
	values = append(values, id)

	addString := func(key string, s *protot.NilString) {
		if s != nil && s.Valid {
			keys = append(keys, key)
			values = append(values, s.String)
		}
	}

	addBool := func(key string, s *protot.NilBool) {
		if s != nil && s.Valid {
			keys = append(keys, key)
			values = append(values, s.Bool)
		}
	}

	addString("username", username)
	addString("password", securePassword)
	addString("first_name", firstName)
	addString("last_name", lastName)

	addBool("is_superuser", isSuperuser)
	addBool("is_active", isActive)
	addBool("is_staff", isStaff)
	addBool("is_confirmed", isConfirmed)

	if len(keys) == 0 {
		return nil, errors.New("charond: nothing to update")
	}

	query := `UPDATE charon.user SET `
	for j, key := range keys {
		if j != 0 {
			query += ", "
		}

		query += fmt.Sprintf("%s = $%d", key, j+2) // plus 2, because of where clause (1+1)
	}

	query += `
		WHERE id = $1
		RETURNING id, password, username, first_name, last_name, is_active, is_staff,
			is_superuser, is_confirmed, confirmation_token, last_login_at,
			created_at, updated_at
	`

	var user userEntity
	err := ur.db.QueryRow(query, values...).Scan(
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

	return &user, nil
}
