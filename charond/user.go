package charond

import (
	"database/sql"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/ntypes"
)

const (
	// UserConfirmationTokenUsed is a value that is used when confirmation token was already used.
	UserConfirmationTokenUsed = "!"
)

// String return concatenated first and last name of the user.
func (ue *userEntity) String() string {
	return ue.FirstName + " " + ue.LastName
}

func (ue *userEntity) message() (*charon.User, error) {
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
		CreatedBy:   ue.CreatedBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   ue.UpdatedBy,
	}, nil
}

type userProvider interface {
	Exists(id int64) (bool, error)
	Create(username string, password []byte, firstName, lastName string, confirmationToken []byte, isSuperuser, isStaff, isActive, isConfirmed bool) (*userEntity, error)
	Insert(*userEntity) (*userEntity, error)
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
		createdBy *ntypes.Int64,
		firstName *ntypes.String,
		isActive *ntypes.Bool,
		isConfirmed *ntypes.Bool,
		isStaff *ntypes.Bool,
		isSuperuser *ntypes.Bool,
		lastLoginAt *time.Time,
		lastName *ntypes.String,
		password []byte,
		updatedAt *time.Time,
		updatedBy *ntypes.Int64,
		username *ntypes.String,
	) (*userEntity, error)
	RegistrationConfirmation(id int64, confirmationToken string) error
	IsGranted(id int64, permission charon.Permission) (bool, error)
	SetPermissions(id int64, permissions ...charon.Permission) (int64, int64, error)
}

func newUserRepository(dbPool *sql.DB) userProvider {
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
	return ur.Insert(entity)
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

func (ur *userRepository) RegistrationConfirmation(userID int64, confirmationToken string) error {
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
func (ur *userRepository) ChangePassword(userID int64, password string) error {
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
func (ur *userRepository) FindOneByUsername(username string) (*userEntity, error) {
	return ur.findOneBy("username", username)
}

func (ur *userRepository) findOneBy(fieldName string, value interface{}) (*userEntity, error) {
	query := `
		SELECT ` + strings.Join(tableUserColumns, ",") + `
		FROM ` + ur.table + `
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

// Exists implements UserRepository interface.
func (ur *userRepository) Exists(userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM ` + ur.table + ` AS p
			WHERE ` + tableUserColumnID + ` = $1
		)
	`

	var exists bool
	if err := ur.db.QueryRow(query, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// IsGranted implements UserRepository interface.
func (ur *userRepository) IsGranted(id int64, p charon.Permission) (bool, error) {
	var exists bool
	subsystem, module, action := p.Split()
	if err := ur.db.QueryRow(isGrantedQuery(
		tableUserPermissions,
		tableUserPermissionsColumnUserID,
		tableUserPermissionsColumnPermissionSubsystem,
		tableUserPermissionsColumnPermissionModule,
		tableUserPermissionsColumnPermissionAction,
	), id, subsystem, module, action).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// SetPermissions implements UserRepository interface.
func (ur *userRepository) SetPermissions(id int64, p ...charon.Permission) (int64, int64, error) {
	return setPermissions(ur.db, tableUserPermissions,
		tableUserPermissionsColumnUserID,
		tableUserPermissionsColumnPermissionSubsystem,
		tableUserPermissionsColumnPermissionModule,
		tableUserPermissionsColumnPermissionAction, id, p)
}

func createDumyTestUser(repo userProvider, hasher charon.PasswordHasher) (*userEntity, error) {
	password, err := hasher.Hash([]byte("test"))
	if err != nil {
		return nil, err
	}
	return repo.CreateSuperuser("test", password, "Test", "Test")
}
