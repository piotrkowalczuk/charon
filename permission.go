package charon

import (
	"bytes"
	"database/sql"
	"errors"
	"strings"
	"sync"
)

const (
	UserCanCreate      Permission = "charon:user:can create"
	UserCanCreateStaff Permission = "charon:user:can create staff"

	UserCanDeleteAsStranger      Permission = "charon:user:can delete as stranger"
	UserCanDeleteAsOwner         Permission = "charon:user:can delete as owner"
	UserCanDeleteStaffAsStranger Permission = "charon:user:can delete staff as stranger"
	UserCanDeleteStaffAsOwner    Permission = "charon:user:can delete staff as owner"

	UserCanModifyAsStranger      Permission = "charon:user:can modify as stranger"
	UserCanModifyAsOwner         Permission = "charon:user:can modify as owner"
	UserCanModifyStaffAsStranger Permission = "charon:user:can modify staff as stranger"
	UserCanModifyStaffAsOwner    Permission = "charon:user:can modify staff as owner"

	UserCanRetrieveAsOwner         Permission = "charon:user:can retrieve as owner"
	UserCanRetrieveAsStranger      Permission = "charon:user:can retrieve as stranger"
	UserCanRetrieveStaffAsOwner    Permission = "charon:user:can retrieve staff as owner"
	UserCanRetrieveStaffAsStranger Permission = "charon:user:can retrieve staff as stranger"

	UserPermissionCanCreate                  Permission = "charon:user_permission:can create"
	UserPermissionCanDelete                  Permission = "charon:user_permission:can delete"
	UserPermissionCanModify                  Permission = "charon:user_permission:can modify"
	UserPermissionCanRetrieve                Permission = "charon:user_permission:can retrieve"
	UserPermissionCanCheckGrantingAsStranger Permission = "charon:user_permission:can check granting as a stranger"

	UserGroupCanCreate                   Permission = "charon:user_group:can create"
	UserGroupCanDelete                   Permission = "charon:user_group:can delete"
	UserGroupCanModify                   Permission = "charon:user_group:can modify"
	UserGroupCanRetrieve                 Permission = "charon:user_group:can retrieve"
	UserGroupCanCheckBelongingAsStranger Permission = "charon:user_group:can check belonging as a stranger"

	PermissionCanCreate   Permission = "charon:permission:can create"
	PermissionCanDelete   Permission = "charon:permission:can delete"
	PermissionCanModify   Permission = "charon:permission:can modify"
	PermissionCanRetrieve Permission = "charon:permission:can retrieve"

	GroupCanCreate   Permission = "charon:group:can create"
	GroupCanDelete   Permission = "charon:group:can delete"
	GroupCanModify   Permission = "charon:group:can modify"
	GroupCanRetrieve Permission = "charon:group:can retrieve"

	GroupPermissionCanCreate   Permission = "charon:group_permission:can create"
	GroupPermissionCanDelete   Permission = "charon:group_permission:can delete"
	GroupPermissionCanModify   Permission = "charon:group_permission:can modify"
	GroupPermissionCanRetrieve Permission = "charon:group_permission:can retrieve"
)

var (
	// EmptyPermission is a shorthand
	EmptyPermission = Permission("")
	// AllPermissions ...
	AllPermissions = Permissions{
		UserCanCreate,
		UserCanCreateStaff,
		UserCanDeleteAsStranger,
		UserCanDeleteAsOwner,
		UserCanDeleteStaffAsStranger,
		UserCanDeleteStaffAsOwner,
		UserCanModifyAsStranger,
		UserCanModifyAsOwner,
		UserCanModifyStaffAsStranger,
		UserCanModifyStaffAsOwner,
		UserCanRetrieveAsOwner,
		UserCanRetrieveAsStranger,
		UserCanRetrieveStaffAsOwner,
		UserCanRetrieveStaffAsStranger,
		UserPermissionCanCreate,
		UserPermissionCanDelete,
		UserPermissionCanModify,
		UserPermissionCanRetrieve,
		PermissionCanCreate,
		PermissionCanDelete,
		PermissionCanModify,
		PermissionCanRetrieve,
		GroupCanCreate,
		GroupCanDelete,
		GroupCanModify,
		GroupCanRetrieve,
		GroupPermissionCanCreate,
		GroupPermissionCanDelete,
		GroupPermissionCanModify,
		GroupPermissionCanRetrieve,
	}
)

// Permission is a string that consist of subsystem, module/content type and an action.
type Permission string

// String implements fmt.Stringer interface.
func (p Permission) String() string {
	return string(p)
}

// Split returns subsystem, module/content ty and action that describes single Permission.
func (p Permission) Split() (string, string, string) {
	if p == "" {
		return "", "", ""
	}

	parts := bytes.Split([]byte(p), []byte(":"))

	switch len(parts) {
	case 1:
		return "", "", string(parts[0])
	case 2:
		return "", string(parts[0]), string(parts[1])
	default:
		return string(parts[0]), string(parts[1]), string(parts[2])
	}
}

// Subsystem is a handy wrapper for Split method, that just returns subsystem.
func (p Permission) Subsystem() (subsystem string) {
	subsystem, _, _ = p.Split()

	return
}

// Module is a handy wrapper for Split method, that just returns module.
func (p Permission) Module() (module string) {
	_, module, _ = p.Split()

	return
}

// Action is a handy wrapper for Split method, that just returns action.
func (p Permission) Action() (action string) {
	_, _, action = p.Split()

	return
}

// Permission implements Permission interface.
func (p Permission) Permission() string {
	return string(p)
}

// MarshalJSON implements json Marshaller interface.
func (p Permission) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return []byte(`""`), nil
	}
	b := make([]byte, 1, len(p))
	b[0] = '"'
	b = append(b, []byte(p)...)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements json Unmarshaler interface.
func (p *Permission) UnmarshalJSON(b []byte) error {
	*p = Permission(string(b))

	return nil
}

type Permissions []Permission

// NewPermissions allocates new Permissions using given slice of strings.
// It maps each string in a slice into Permission.
func NewPermissions(ss ...string) Permissions {
	ps := make(Permissions, 0, len(ss))
	for _, s := range ss {
		ps = append(ps, Permission(s))
	}

	return ps
}

// Contains returns true if given Permission exists in the collection.
// If none is provided returns false.
func (p Permissions) Contains(permissions ...Permission) bool {
	if len(permissions) == 0 {
		return false
	}

	for _, perm := range p {
		for _, pp := range permissions {
			if perm == pp {
				return true
			}
		}
	}

	return false
}

// Strings maps Permissions into slice of strings.
func (p Permissions) Strings() (s []string) {
	s = make([]string, 0, len(p))
	for _, pp := range p {
		s = append(s, pp.String())
	}

	return s
}

// Permission returns Permission value that is concatenated
// using entity properties like subsystem, module and action.
func (pe *permissionEntity) Permission() Permission {
	return Permission(pe.Subsystem + ":" + pe.Module + ":" + pe.Action)
}

type permissionProvider interface {
	Find(criteria *permissionCriteria) ([]*permissionEntity, error)
	FindOneByID(id int64) (entity *permissionEntity, err error)
	// FindByUserID retrieves all permissions for user represented by given id.
	FindByUserID(userID int64) (entities []*permissionEntity, err error)
	// FindByGroupID retrieves all permissions for group represented by given id.
	FindByGroupID(groupID int64) (entities []*permissionEntity, err error)
	Register(permissions Permissions) (created, untouched, removed int64, err error)
	Insert(entity *permissionEntity) (*permissionEntity, error)
}

func newPermissionRepository(dbPool *sql.DB) *permissionRepository {
	return &permissionRepository{
		db:      dbPool,
		table:   tablePermission,
		columns: tablePermissionColumns,
	}
}

// FindByUserID implements PermissionRepository interface.
func (pr *permissionRepository) FindByUserID(userID int64) ([]*permissionEntity, error) {
	// TODO: does it work?
	return pr.findBy(`
		SELECT DISTINCT ON (p.id)
			`+columns(tablePermissionColumns, "p")+`
		FROM `+pr.table+` AS p
		LEFT JOIN `+tableUserPermissions+` AS up
			ON up.`+tableUserPermissionsColumnPermissionSubsystem+` = p.`+tablePermissionColumnSubsystem+`
			AND up.`+tableUserPermissionsColumnPermissionModule+` = p.`+tablePermissionColumnModule+`
			AND up.`+tableUserPermissionsColumnPermissionAction+` = p.`+tablePermissionColumnAction+`
		LEFT JOIN `+tableUserGroups+` AS ug ON ug.`+tableUserGroupsColumnUserID+` = $1
		LEFT JOIN `+tableGroupPermissions+` AS gp
			ON gp.`+tableGroupPermissionsColumnPermissionSubsystem+` = p.`+tablePermissionColumnSubsystem+`
			AND gp.`+tableGroupPermissionsColumnPermissionModule+` = p.`+tablePermissionColumnModule+`
			AND gp.`+tableGroupPermissionsColumnPermissionAction+` = p.`+tablePermissionColumnAction+`
			AND gp.`+tableGroupPermissionsColumnGroupID+` = ug.`+tableUserGroupsColumnGroupID+`
		WHERE up.`+tableUserPermissionsColumnUserID+` = $1 OR ug.`+tableUserGroupsColumnUserID+` = $1
	`, userID)
}

// FindByGroupID implements PermissionRepository interface.
func (pr *permissionRepository) FindByGroupID(userID int64) ([]*permissionEntity, error) {
	// TODO: does it work?
	return pr.findBy(`
		SELECT DISTINCT ON (p.id)
			`+columns(tablePermissionColumns, "p")+`
		FROM `+pr.table+` AS p
		LEFT JOIN `+tableGroupPermissions+` AS gp ON gp.permission_id = p.id AND gp.group_id = ug.group_id
		WHERE up.user_id = $1 OR ug.user_id = $1
	`, userID)
}

func (pr *permissionRepository) findBy(query string, args ...interface{}) ([]*permissionEntity, error) {
	rows, err := pr.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []*permissionEntity{}
	for rows.Next() {
		var p permissionEntity
		err = rows.Scan(
			&p.Action,
			&p.CreatedAt,
			&p.ID,
			&p.Module,
			&p.Subsystem,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, &p)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return permissions, nil
}

func (pr *permissionRepository) findOneStmt() (*sql.Stmt, error) {
	return pr.db.Prepare(
		"SELECT " + strings.Join(tablePermissionColumns, ",") + " " +
			"FROM " + pr.table + " AS p " +
			"WHERE p.subsystem = $1 AND p.module = $2 AND p.action = $3",
	)
}

// Register implements PermissionRepository interface.
func (pr *permissionRepository) Register(permissions Permissions) (created, unt, removed int64, err error) {
	var (
		tx             *sql.Tx
		insert, delete *sql.Stmt
		rows           *sql.Rows
		res            sql.Result
		subsystem      string
		entities       []*permissionEntity
		affected       int64
	)
	if len(permissions) == 0 {
		return 0, 0, 0, errors.New("charon: empty slice, permissions cannot be registered")
	}

	subsystem = permissions[0].Subsystem()
	if subsystem == "" {
		return 0, 0, 0, errors.New("charon: subsystem name is empty string, permissions cannot be registered")
	}

	for _, p := range permissions {
		if p.Subsystem() != subsystem {
			return 0, 0, 0, errors.New("charon: provided permissions do not belong to one subsystem, permissions cannot be registered")
		}
	}

	tx, err = pr.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			unt = untouched(int64(len(permissions)), created, removed)
		}
	}()

	rows, err = tx.Query("SELECT "+strings.Join(tablePermissionColumns, ",")+" FROM "+pr.table+" AS p WHERE p.subsystem = $1", subsystem)
	if err != nil {
		return
	}
	defer rows.Close()

	entities = []*permissionEntity{}
	for rows.Next() {
		var entity permissionEntity
		err = rows.Scan(
			&entity.Action,
			&entity.CreatedAt,
			&entity.ID,
			&entity.Module,
			&entity.Subsystem,
			&entity.UpdatedAt,
		)
		if err != nil {
			return
		}
		entities = append(entities, &entity)
	}
	if rows.Err() != nil {
		return 0, 0, 0, rows.Err()
	}

	insert, err = tx.Prepare("INSERT INTO " + pr.table + " (subsystem, module, action) VALUES ($1, $2, $3)")
	if err != nil {
		return
	}

MissingPermissionsLoop:
	for _, p := range permissions {
		for _, e := range entities {
			if p == e.Permission() {
				continue MissingPermissionsLoop
			}
		}

		if res, err = insert.Exec(p.Split()); err != nil {
			return
		}
		if affected, err = res.RowsAffected(); err != nil {
			return
		}
		created += affected
	}

	delete, err = tx.Prepare("DELETE FROM " + pr.table + " AS p WHERE p.id = $1")
	if err != nil {
		return
	}

RedundantPermissionsLoop:
	for _, e := range entities {
		for _, p := range permissions {
			if e.Permission() == p {
				continue RedundantPermissionsLoop
			}
		}

		if res, err = delete.Exec(e.ID); err != nil {
			return
		}
		if affected, err = res.RowsAffected(); err != nil {
			return
		}

		removed += affected
	}

	return
}

// PermissionRegistry is an interface that describes in memory storage that holds information
// about permissions that was registered by 3rd party services.
// Should be only used as a proxy for registration process to avoid multiple sql hits.
type PermissionRegistry interface {
	// Exists returns true if given Permission was already registered.
	Exists(permission Permission) (exists bool)
	// Register checks if given collection is valid and
	// calls PermissionRepository to store provided permissions
	// in persistent way.
	Register(permissions Permissions) (created, untouched, removed int64, err error)
}

type permissionRegistry struct {
	sync.RWMutex
	repository  permissionProvider
	permissions map[Permission]struct{}
}

func newPermissionRegistry(r permissionProvider) PermissionRegistry {
	return &permissionRegistry{
		repository:  r,
		permissions: make(map[Permission]struct{}),
	}
}

// Exists implements PermissionRegistry interface.
func (pr *permissionRegistry) Exists(permission Permission) (ok bool) {
	pr.RLock()
	pr.RUnlock()

	_, ok = pr.permissions[permission]
	return
}

// Register implements PermissionRegistry interface.
func (pr *permissionRegistry) Register(permissions Permissions) (created, untouched, removed int64, err error) {
	pr.Lock()
	defer pr.Unlock()

	nb := 0
	for _, p := range permissions {
		if _, ok := pr.permissions[p]; !ok {
			pr.permissions[p] = struct{}{}
			nb++
		}
	}

	if nb > 0 {
		return pr.repository.Register(permissions)
	}

	return 0, 0, 0, nil
}

// FindByTag implements PermissionRepository interface.
func (pr *permissionRepository) FindByTag(userID int64) ([]*permissionEntity, error) {
	query := `
		SELECT DISTINCT ON (p.id)
			` + columns(tablePermissionColumns, "p") + `
		FROM ` + pr.table + ` AS p
		LEFT JOIN ` + tableUserPermissions + ` AS up ON up.permission_id = p.id AND up.user_id = $1
		LEFT JOIN ` + tableUserGroups + ` AS ug ON ug.user_id = $1
		LEFT JOIN ` + tableGroupPermissions + ` AS gp ON gp.permission_id = p.id AND gp.group_id = ug.group_id
		WHERE up.user_id = $1 OR ug.user_id = $1
	`

	rows, err := pr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []*permissionEntity{}
	for rows.Next() {
		var p permissionEntity
		err = rows.Scan(
			&p.Action,
			&p.CreatedAt,
			&p.ID,
			&p.Module,
			&p.Subsystem,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, &p)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return permissions, nil
}
