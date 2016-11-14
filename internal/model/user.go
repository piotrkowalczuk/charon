package model

import (
	"database/sql"
	"strings"

	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
)

const (
	// UserConfirmationTokenUsed ...
	UserConfirmationTokenUsed = "!"
)

var (
	// ExternalPassword is a password that is set when external source of authentication is provided (e.g. LDAP).
	ExternalPassword = []byte("!")
)

// String return concatenated first and last name of the user.
// Implements fmt Stringer interface.
func (ue *UserEntity) String() string {
	switch {
	case ue.FirstName == "":
		return ue.LastName
	case ue.LastName == "":
		return ue.FirstName
	default:
		return ue.FirstName + " " + ue.LastName
	}
}

// Message maps entity into protobuf message.
func (ue *UserEntity) Message() (*charonrpc.User, error) {
	var (
		err                  error
		createdAt, updatedAt *pbts.Timestamp
	)

	if createdAt, err = ptypes.TimestampProto(ue.CreatedAt); err != nil {
		return nil, err
	}
	if ue.UpdatedAt != nil {
		if createdAt, err = ptypes.TimestampProto(*ue.UpdatedAt); err != nil {
			return nil, err
		}
	}

	return &charonrpc.User{
		Id:          ue.ID,
		Username:    ue.Username,
		FirstName:   ue.FirstName,
		LastName:    ue.LastName,
		IsSuperuser: ue.IsSuperuser,
		IsActive:    ue.IsActive,
		IsStaff:     ue.IsStaff,
		IsConfirmed: ue.IsConfirmed,
		CreatedAt:   createdAt,
		CreatedBy:   ue.CreatedBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   ue.UpdatedBy,
	}, nil
}

// UserProvider wraps UserRepository into interface.
type UserProvider interface {
	Exists(id int64) (bool, error)
	Create(username string, password []byte, FirstName, LastName string, confirmationToken []byte, isSuperuser, IsStaff, isActive, isConfirmed bool) (*UserEntity, error)
	Insert(*UserEntity) (*UserEntity, error)
	CreateSuperuser(username string, password []byte, FirstName, LastName string) (*UserEntity, error)
	// Count retrieves number of all users.
	Count() (int64, error)
	UpdateLastLoginAt(id int64) (int64, error)
	ChangePassword(id int64, password string) error
	Find(criteria *UserCriteria) ([]*UserEntity, error)
	FindOneByID(id int64) (*UserEntity, error)
	FindOneByUsername(username string) (*UserEntity, error)
	DeleteOneByID(id int64) (int64, error)
	UpdateOneByID(int64, *UserPatch) (*UserEntity, error)
	RegistrationConfirmation(id int64, confirmationToken string) error
	IsGranted(id int64, permission charon.Permission) (bool, error)
	SetPermissions(id int64, permissions ...charon.Permission) (int64, int64, error)
}

// UserRepository extends UserRepositoryBase.
type UserRepository struct {
	UserRepositoryBase
}

// NewUserRepository alocates new UserRepository instance
func NewUserRepository(dbPool *sql.DB) *UserRepository {
	return &UserRepository{
		UserRepositoryBase: UserRepositoryBase{
			db:      dbPool,
			table:   TableUser,
			columns: TableUserColumns,
		},
	}
}

// Create implements UserProvider interface.
func (ur *UserRepository) Create(username string, password []byte, firstName, lastName string, confirmationToken []byte, isSuperuser, isStaff, isActive, isConfirmed bool) (*UserEntity, error) {
	if isSuperuser {
		isStaff = true
		isActive = true
		isConfirmed = true
	}

	entity := &UserEntity{
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
	return ur.Insert(entity)
}

// CreateSuperuser implements UserProvider interface.
func (ur *UserRepository) CreateSuperuser(username string, password []byte, FirstName, LastName string) (*UserEntity, error) {
	return ur.Create(username, password, FirstName, LastName, []byte(UserConfirmationTokenUsed), true, false, true, true)
}

// Count implements UserProvider interface.
func (ur *UserRepository) Count() (n int64, err error) {
	err = ur.db.QueryRow("SELECT COUNT(*) FROM " + TableUser).Scan(&n)

	return
}

// RegistrationConfirmation ...
func (ur *UserRepository) RegistrationConfirmation(userID int64, confirmationToken string) error {
	query := `
		UPDATE ` + ur.table + `
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
func (ur *UserRepository) ChangePassword(userID int64, password string) error {
	query := `
		UPDATE ` + ur.table + `
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
func (ur *UserRepository) FindOneByUsername(username string) (*UserEntity, error) {
	return ur.FindOneBy("username", username)
}

// FindOneBy ...
func (ur *UserRepository) FindOneBy(fieldName string, value interface{}) (*UserEntity, error) {
	query := `
		SELECT ` + strings.Join(TableUserColumns, ",") + `
		FROM ` + ur.table + `
		WHERE ` + fieldName + ` = $1
		LIMIT 1
	`

	var ent UserEntity
	err := ur.db.QueryRow(query, value).Scan(
		&ent.ConfirmationToken,
		&ent.CreatedAt,
		&ent.CreatedBy,
		&ent.FirstName,
		&ent.ID,
		&ent.IsActive,
		&ent.IsConfirmed,
		&ent.IsStaff,
		&ent.IsSuperuser,
		&ent.LastLoginAt,
		&ent.LastName,
		&ent.Password,
		&ent.UpdatedAt,
		&ent.UpdatedBy,
		&ent.Username,
	)

	if err != nil {
		return nil, err
	}

	return &ent, nil
}

// UpdateLastLoginAt implements UserProvider interface.
func (ur *UserRepository) UpdateLastLoginAt(userID int64) (int64, error) {
	query := `
		UPDATE ` + ur.table + `
		SET ` + TableUserColumnLastLoginAt + ` = NOW()
		WHERE ` + TableUserColumnID + ` = $1
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

// Exists implements UserProvider interface.
func (ur *UserRepository) Exists(userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM ` + ur.table + ` AS p
			WHERE ` + TableUserColumnID + ` = $1
		)
	`

	var exists bool
	if err := ur.db.QueryRow(query, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// IsGranted implements UserProvider interface.
func (ur *UserRepository) IsGranted(id int64, p charon.Permission) (bool, error) {
	var exists bool
	subsystem, module, action := p.Split()
	if err := ur.db.QueryRow(isGrantedQuery(
		TableUserPermissions,
		TableUserPermissionsColumnUserID,
		TableUserPermissionsColumnPermissionSubsystem,
		TableUserPermissionsColumnPermissionModule,
		TableUserPermissionsColumnPermissionAction,
	), id, subsystem, module, action).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// SetPermissions implements UserProvider interface.
func (ur *UserRepository) SetPermissions(id int64, p ...charon.Permission) (int64, int64, error) {
	return setPermissions(ur.db, TableUserPermissions,
		TableUserPermissionsColumnUserID,
		TableUserPermissionsColumnPermissionSubsystem,
		TableUserPermissionsColumnPermissionModule,
		TableUserPermissionsColumnPermissionAction, id, p)
}
