package main

import (
	"database/sql"
	"strings"

	"sync"

	"errors"

	"github.com/piotrkowalczuk/charon"
)

// Permission returns Permission value that is concatenated
// using entity properties like subsystem, module and action.
func (pe *permissionEntity) Permission() charon.Permission {
	return charon.Permission(pe.Subsystem + ":" + pe.Module + ":" + pe.Action)
}

// PermissionRepository ...
type PermissionRepository interface {
	Find(criteria *permissionCriteria) ([]*permissionEntity, error)
	FindOneByID(id int64) (entity *permissionEntity, err error)
	FindByUserID(userID int64) (entities []*permissionEntity, err error)
	Register(permissions charon.Permissions) (created, untouched, removed int, err error)
	Insert(entity *permissionEntity) (*permissionEntity, error)
}

func newPermissionRepository(dbPool *sql.DB) *permissionRepository {
	return &permissionRepository{
		db:      dbPool,
		table:   tablePermission,
		columns: tablePermissionColumns,
	}
}

// FindByUserID retrieves all permissions for user represented by given id.
func (pr *permissionRepository) FindByUserID(userID int64) ([]*permissionEntity, error) {
	query := `
		SELECT DISTINCT ON (p.id)
			` + columns(tablePermissionColumns, "p") + `
		FROM ` + tablePermission + ` AS p
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

func (pr *permissionRepository) findOneStmt() (*sql.Stmt, error) {
	return pr.db.Prepare(
		"SELECT " + strings.Join(tablePermissionColumns, ",") + " " +
			"FROM " + tablePermission + " AS p " +
			"WHERE p.subsystem = $1 AND p.module = $2 AND p.action = $3",
	)
}

func (pr *permissionRepository) Register(permissions charon.Permissions) (created, untouched, removed int, err error) {
	var (
		tx             *sql.Tx
		insert, delete *sql.Stmt
		rows           *sql.Rows
		subsystem      string
		entities       []*permissionEntity
	)
	if len(permissions) == 0 {
		return 0, 0, 0, errors.New("charond: empty slice, permissions cannot be registered")
	}

	subsystem = permissions[0].Subsystem()
	if subsystem == "" {
		return 0, 0, 0, errors.New("charond: subsystem name is empty string, permissions cannot be registered")
	}

	for _, p := range permissions {
		if p.Subsystem() != subsystem {
			return 0, 0, 0, errors.New("charond: provided permissions do not belong to one subsystem, permissions cannot be registered")
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
			tx.Commit()
			untouched = len(permissions) - removed - created
		}
	}()

	rows, err = tx.Query("SELECT "+strings.Join(tablePermissionColumns, ",")+" FROM "+tablePermission+" AS p WHERE p.subsystem = $1", subsystem)
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

	insert, err = tx.Prepare("INSERT INTO " + tablePermission + " (subsystem, module, action) VALUES ($1, $2, $3)")
	if err != nil {
		return
	}

MissingPermissionsLoop:
	for _, p := range permissions {
		for _, e := range entities {
			if p == e.Permission() {
				break MissingPermissionsLoop
			}
		}

		if _, err = insert.Exec(p.Split()); err != nil {
			return
		}
		created++
	}

	delete, err = tx.Prepare("DELETE FROM " + tablePermission + " AS p WHERE p.id = $1")
	if err != nil {
		return
	}

RedundantPermissionsLoop:
	for _, e := range entities {
		for _, p := range permissions {
			if e.Permission() == p {
				break RedundantPermissionsLoop
			}
		}

		if _, err = delete.Exec(e.ID); err != nil {
			return
		}
		removed++
	}

	return
}

// PermissionRegistry is an interface that describes in memory storage that holds information
// about permissions that was registered by 3rd party services.
// Should be only used as a proxy for registration process to avoid multiple sql hits.
type PermissionRegistry interface {
	// Exists returns true if given Permission was already registered.
	Exists(permission charon.Permission) (exists bool)
	// Register checks if given collection is valid and
	// calls PermissionRepository to store provided permissions
	// in persistent way.
	Register(permissions charon.Permissions) (created, untouched, removed int, err error)
}

type permissionRegistry struct {
	sync.RWMutex
	repository  PermissionRepository
	permissions map[charon.Permission]struct{}
}

func newPermissionRegistry(r PermissionRepository) PermissionRegistry {
	return &permissionRegistry{
		repository:  r,
		permissions: make(map[charon.Permission]struct{}),
	}
}

// Exists implements PermissionRegistry interface.
func (pr *permissionRegistry) Exists(permission charon.Permission) (ok bool) {
	pr.RLock()
	pr.RUnlock()

	_, ok = pr.permissions[permission]
	return
}

// Register implements PermissionRegistry interface.
func (pr *permissionRegistry) Register(permissions charon.Permissions) (created, untouched, removed int, err error) {
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
