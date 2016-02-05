package main

import (
	"database/sql"
	"strings"
	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
	"github.com/piotrkowalczuk/protot"
)

const (
	// UserConfirmationTokenUsed is a value that is used when confirmation token was already used.
	UserConfirmationTokenUsed = "!"
)

// String return concatenated first and last name of the user.
func (ue *userEntity) String() string {
	return ue.FirstName + " " + ue.LastName
}

// Message allocates new corresponding protobuf message.
func (ue *userEntity) Message() *charon.User {
	var (
		createdAt *protot.Timestamp
		updatedAt *protot.Timestamp
	)

	createdAt = protot.TimeToTimestamp(ue.CreatedAt)
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
		CreatedBy:   &ue.CreatedBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   &ue.UpdatedBy,
	}
}

// UserRepository ...
type UserRepository interface {
	Create(username string, password []byte, firstName, lastName string, confirmationToken []byte, isSuperuser, isStaff, isActive, isConfirmed bool) (*userEntity, error)
	CreateSuperuser(username string, password []byte, firstName, lastName string) (*userEntity, error)
	// Count retrieves number of all users.
	Count() (int64, error)
	UpdateLastLoginAt(id int64) (int64, error)
	ChangePassword(id int64, password string) error
	Find(criteria *userCriteria) ([]*userEntity, error)
	FindOneByID(id int64) (*userEntity, error)
	FindOneByUsername(username string) (*userEntity, error)
	DeleteByID(id int64) (int64, error)
	UpdateByID(
		id int64,
		confirmationToken []byte,
		createdAt *time.Time,
		createdBy nilt.Int64,
		firstName nilt.String,
		isActive nilt.Bool,
		isConfirmed nilt.Bool,
		isStaff nilt.Bool,
		isSuperuser nilt.Bool,
		lastLoginAt *time.Time,
		lastName nilt.String,
		password []byte,
		updatedAt *time.Time,
		updatedBy nilt.Int64,
		username nilt.String,
	) (*userEntity, error)
	RegistrationConfirmation(id int64, confirmationToken string) error
}

func newUserRepository(dbPool *sql.DB) UserRepository {
	return &userRepository{
		db:      dbPool,
		table:   tableUser,
		columns: tableUserColumns,
	}
}

// Create implements UserRepository interface.
func (ur *userRepository) Create(username string, password []byte, firstName, lastName string, confirmationToken []byte, isSuperuser, isStaff, isActive, isConfirmed bool) (*userEntity, error) {
	if isSuperuser {
		isStaff = true
		isActive = true
		isConfirmed = true
	}

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
func (ur *userRepository) CreateSuperuser(username string, password []byte, firstName, lastName string) (*userEntity, error) {
	return ur.Create(username, password, firstName, lastName, []byte(UserConfirmationTokenUsed), true, false, true, true)
}

// Count implements UserRepository interface.
func (ur *userRepository) Count() (n int64, err error) {
	err = ur.db.QueryRow("SELECT COUNT(*) FROM " + tableUser).Scan(&n)

	return
}

func (ur *userRepository) insert(e *userEntity) error {
	query := `
		INSERT INTO ` + tableUser + ` (
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
		UPDATE ` + tableUser + `
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
		UPDATE ` + tableUser + `
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

func (ur *userRepository) findOneBy(fieldName string, value interface{}) (*userEntity, error) {
	query := `
		SELECT ` + strings.Join(tableUserColumns, ",") + `
		FROM ` + tableUser + `
		WHERE ` + fieldName + ` = $1
		LIMIT 1
	`

	var user userEntity
	err := ur.db.QueryRow(query, value).Scan(
		&user.ConfirmationToken,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.FirstName,
		&user.ID,
		&user.IsActive,
		&user.IsConfirmed,
		&user.IsStaff,
		&user.IsSuperuser,
		&user.LastLoginAt,
		&user.LastName,
		&user.Password,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.Username,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateLastLoginAt implements UserRepository interface.
func (ur *userRepository) UpdateLastLoginAt(userID int64) (int64, error) {
	query := `
		UPDATE ` + ur.table + `
		SET ` + tableUserColumnLastLoginAt + ` = NOW()
		WHERE ` + tableUserColumnID + ` = $1
	`

	result, err := ur.db.Exec(
		query,
		userID,
	)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
