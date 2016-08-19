package charond

import (
	"database/sql"
	"strings"

	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/piotrkowalczuk/charon"
)

const (
	// UserConfirmationTokenUsed is a value that is used when confirmation token was already used.
	UserConfirmationTokenUsed = "!"
)

// String return concatenated first and last name of the user.
// Implements fmt Stringer interface.
func (ue *userEntity) String() string {
	switch {
	case ue.firstName == "":
		return ue.lastName
	case ue.lastName == "":
		return ue.firstName
	default:
		return ue.firstName + " " + ue.lastName
	}
}

func (ue *userEntity) message() (*charon.User, error) {
	var (
		err                  error
		createdAt, updatedAt *pbts.Timestamp
	)

	if createdAt, err = ptypes.TimestampProto(ue.createdAt); err != nil {
		return nil, err
	}
	if ue.updatedAt != nil {
		if createdAt, err = ptypes.TimestampProto(*ue.updatedAt); err != nil {
			return nil, err
		}
	}

	return &charon.User{
		Id:          ue.id,
		Username:    ue.username,
		FirstName:   ue.firstName,
		LastName:    ue.lastName,
		IsSuperuser: ue.isSuperuser,
		IsActive:    ue.isActive,
		IsStaff:     ue.isStaff,
		IsConfirmed: ue.isConfirmed,
		CreatedAt:   createdAt,
		CreatedBy:   ue.createdBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   ue.updatedBy,
	}, nil
}

type userProvider interface {
	Exists(id int64) (bool, error)
	Create(username string, password []byte, firstName, lastName string, confirmationToken []byte, isSuperuser, isStaff, isActive, isConfirmed bool) (*userEntity, error)
	insert(*userEntity) (*userEntity, error)
	CreateSuperuser(username string, password []byte, firstName, lastName string) (*userEntity, error)
	// Count retrieves number of all users.
	Count() (int64, error)
	updateLastLoginAt(id int64) (int64, error)
	ChangePassword(id int64, password string) error
	find(criteria *userCriteria) ([]*userEntity, error)
	findOneByID(id int64) (*userEntity, error)
	findOneByUsername(username string) (*userEntity, error)
	deleteOneByID(id int64) (int64, error)
	updateOneByID(int64, *userPatch) (*userEntity, error)
	RegistrationConfirmation(id int64, confirmationToken string) error
	IsGranted(id int64, permission charon.Permission) (bool, error)
	SetPermissions(id int64, permissions ...charon.Permission) (int64, int64, error)
}

type userRepository struct {
	userRepositoryBase
}

func newUserRepository(dbPool *sql.DB) userProvider {
	return &userRepository{
		userRepositoryBase: userRepositoryBase{
			db:      dbPool,
			table:   tableUser,
			columns: tableUserColumns,
		},
	}
}

// Create implements userProvider interface.
func (ur *userRepository) Create(username string, password []byte, firstName, lastName string, confirmationToken []byte, isSuperuser, isStaff, isActive, isConfirmed bool) (*userEntity, error) {
	if isSuperuser {
		isStaff = true
		isActive = true
		isConfirmed = true
	}

	entity := &userEntity{
		username:          username,
		password:          password,
		firstName:         firstName,
		lastName:          lastName,
		confirmationToken: confirmationToken,
		isSuperuser:       isSuperuser,
		isStaff:           isStaff,
		isActive:          isActive,
		isConfirmed:       isConfirmed,
	}
	return ur.insert(entity)
}

// CreateSuperuser implements userProvider interface.
func (ur *userRepository) CreateSuperuser(username string, password []byte, firstName, lastName string) (*userEntity, error) {
	return ur.Create(username, password, firstName, lastName, []byte(UserConfirmationTokenUsed), true, false, true, true)
}

// Count implements userProvider interface.
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

// findOneByUsername ...
func (ur *userRepository) findOneByUsername(username string) (*userEntity, error) {
	return ur.findOneBy("username", username)
}

func (ur *userRepository) findOneBy(fieldName string, value interface{}) (*userEntity, error) {
	query := `
		SELECT ` + strings.Join(tableUserColumns, ",") + `
		FROM ` + ur.table + `
		WHERE ` + fieldName + ` = $1
		LIMIT 1
	`

	var ent userEntity
	err := ur.db.QueryRow(query, value).Scan(
		&ent.confirmationToken,
		&ent.createdAt,
		&ent.createdBy,
		&ent.firstName,
		&ent.id,
		&ent.isActive,
		&ent.isConfirmed,
		&ent.isStaff,
		&ent.isSuperuser,
		&ent.lastLoginAt,
		&ent.lastName,
		&ent.password,
		&ent.updatedAt,
		&ent.updatedBy,
		&ent.username,
	)

	if err != nil {
		return nil, err
	}

	return &ent, nil
}

// updateLastLoginAt implements userProvider interface.
func (ur *userRepository) updateLastLoginAt(userID int64) (int64, error) {
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

// Exists implements userProvider interface.
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

// IsGranted implements userProvider interface.
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

// SetPermissions implements userProvider interface.
func (ur *userRepository) SetPermissions(id int64, p ...charon.Permission) (int64, int64, error) {
	return setPermissions(ur.db, tableUserPermissions,
		tableUserPermissionsColumnUserID,
		tableUserPermissionsColumnPermissionSubsystem,
		tableUserPermissionsColumnPermissionModule,
		tableUserPermissionsColumnPermissionAction, id, p)
}

func createDummyTestUser(repo userProvider, hasher charon.PasswordHasher) (*userEntity, error) {
	password, err := hasher.Hash([]byte("test"))
	if err != nil {
		return nil, err
	}
	return repo.CreateSuperuser("test", password, "Test", "Test")
}
