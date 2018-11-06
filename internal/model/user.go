package model

import (
	"context"
	"database/sql"
	"strings"

	"github.com/piotrkowalczuk/charon"
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

// UserProvider wraps UserRepository into interface.
type UserProvider interface {
	Exists(context.Context, int64) (bool, error)
	Create(context.Context, *UserEntity) (*UserEntity, error)
	Insert(context.Context, *UserEntity) (*UserEntity, error)
	CreateSuperuser(ctx context.Context, username string, password []byte, FirstName, LastName string) (*UserEntity, error)
	// Count retrieves number of all users.
	Count(context.Context) (int64, error)
	UpdateLastLoginAt(ctx context.Context, id int64) (int64, error)
	ChangePassword(ctx context.Context, id int64, password string) error
	Find(context.Context, *UserFindExpr) ([]*UserEntity, error)
	FindOneByID(context.Context, int64) (*UserEntity, error)
	FindOneByUsername(context.Context, string) (*UserEntity, error)
	DeleteOneByID(context.Context, int64) (int64, error)
	UpdateOneByID(context.Context, int64, *UserPatch) (*UserEntity, error)
	RegistrationConfirmation(ctx context.Context, id int64, confirmationToken string) (int64, error)
	IsGranted(ctx context.Context, id int64, permission charon.Permission) (bool, error)
	SetPermissions(ctx context.Context, id int64, permissions ...charon.Permission) (int64, int64, error)
}

// UserRepository extends UserRepositoryBase.
type UserRepository struct {
	UserRepositoryBase
}

// NewUserRepository alocates new UserRepository instance
func NewUserRepository(dbPool *sql.DB) *UserRepository {
	return &UserRepository{
		UserRepositoryBase: UserRepositoryBase{
			DB:      dbPool,
			Table:   TableUser,
			Columns: TableUserColumns,
		},
	}
}

// Create implements UserProvider interface.
func (ur *UserRepository) Create(ctx context.Context, ent *UserEntity) (*UserEntity, error) {
	tmp := *ent
	if tmp.IsSuperuser {
		tmp.IsStaff = false
		tmp.IsActive = true
		tmp.IsConfirmed = true
	}

	return ur.Insert(ctx, &tmp)
}

// CreateSuperuser implements UserProvider interface.
func (ur *UserRepository) CreateSuperuser(ctx context.Context, username string, password []byte, firstName, lastName string) (*UserEntity, error) {
	return ur.Create(ctx, &UserEntity{
		Username:          username,
		Password:          password,
		FirstName:         firstName,
		LastName:          lastName,
		ConfirmationToken: []byte(UserConfirmationTokenUsed),
		IsSuperuser:       true,
	})
}

// Count implements UserProvider interface.
func (ur *UserRepository) Count(ctx context.Context) (n int64, err error) {
	err = ur.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+TableUser).Scan(&n)

	return
}

// RegistrationConfirmation ...
func (ur *UserRepository) RegistrationConfirmation(ctx context.Context, userID int64, confirmationToken string) (int64, error) {
	query := `
		UPDATE ` + ur.Table + `
		SET is_confirmed = true, is_active = true, updated_at = NOW(), confirmation_token = $1
		WHERE is_confirmed = false AND id = $2 AND confirmation_token = $3;
	`

	result, err := ur.DB.ExecContext(ctx, query, UserConfirmationTokenUsed, userID, confirmationToken)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// ChangePassword ...
func (ur *UserRepository) ChangePassword(ctx context.Context, userID int64, password string) error {
	query := `
		UPDATE ` + ur.Table + `
		SET password = $2, updated_at = NOW()
		WHERE id = $1;
	`

	result, err := ur.DB.ExecContext(ctx, query, userID, password)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()

	return err
}

// FindOneByUsername ...
func (ur *UserRepository) FindOneByUsername(ctx context.Context, username string) (*UserEntity, error) {
	return ur.FindOneBy(ctx, "username", username)
}

// FindOneBy ...
func (ur *UserRepository) FindOneBy(ctx context.Context, fieldName string, value interface{}) (*UserEntity, error) {
	query := `
		SELECT ` + strings.Join(TableUserColumns, ",") + `
		FROM ` + ur.Table + `
		WHERE ` + fieldName + ` = $1
		LIMIT 1
	`

	var ent UserEntity
	err := ur.DB.QueryRowContext(ctx, query, value).Scan(
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
func (ur *UserRepository) UpdateLastLoginAt(ctx context.Context, userID int64) (int64, error) {
	query := `
		UPDATE ` + ur.Table + `
		SET ` + TableUserColumnLastLoginAt + ` = NOW()
		WHERE ` + TableUserColumnID + ` = $1
	`

	result, err := ur.DB.ExecContext(ctx,
		query,
		userID,
	)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Exists implements UserProvider interface.
func (ur *UserRepository) Exists(ctx context.Context, userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM ` + ur.Table + ` AS p
			WHERE ` + TableUserColumnID + ` = $1
		)
	`

	var exists bool
	if err := ur.DB.QueryRowContext(ctx, query, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// IsGranted implements UserProvider interface.
func (ur *UserRepository) IsGranted(ctx context.Context, id int64, p charon.Permission) (bool, error) {
	var exists bool
	subsystem, module, action := p.Split()
	if err := ur.DB.QueryRowContext(ctx, isGrantedQuery(
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
func (ur *UserRepository) SetPermissions(ctx context.Context, id int64, p ...charon.Permission) (int64, int64, error) {
	return setPermissions(ur.DB, ctx, TableUserPermissions,
		TableUserPermissionsColumnUserID,
		TableUserPermissionsColumnPermissionSubsystem,
		TableUserPermissionsColumnPermissionModule,
		TableUserPermissionsColumnPermissionAction, id, p)
}
