package main

import (
	"database/sql"
	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/pqcnstr"
)

const (
	tablePermission                                                         = "charon.permission"
	tablePermissionConstraintPrimaryKey                  pqcnstr.Constraint = tablePermission + "_pkey"
	tablePermissionConstraintUniqueSubsystemModuleAction pqcnstr.Constraint = tablePermission + "_subsystem_module_action_key"
	tablePermissionConstraintForeignKeyUpdatedBy         pqcnstr.Constraint = tablePermission + "_created_by_fkey"
	tablePermissionCreate                                                   = `
		CREATE TABLE IF NOT EXISTS ` + tablePermission + ` (
			id           SERIAL,
			subsystem    TEXT                      NOT NULL,
			module       TEXT                      NOT NULL,
			action       TEXT                      NOT NULL,
			created_at   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
			created_by   INTEGER,

			CONSTRAINT "` + tablePermissionConstraintPrimaryKey + `" PRIMARY KEY (id),
			CONSTRAINT "` + tablePermissionConstraintUniqueSubsystemModuleAction + `" UNIQUE (subsystem, module, action),
			CONSTRAINT "` + tablePermissionConstraintForeignKeyUpdatedBy + `" FOREIGN KEY (created_by) REFERENCES ` + tableUser + ` (id)
		)
	`
	tablePermissionColumns = `
		p.id, p.subsystem, p.module, p.action, p.created_at, p.created_by
	`
)

type permissionEntity struct {
	ID          int64
	SubsystemID int64
	Subsystem   string
	Module      string
	Action      string
	CreatedAt   *time.Time
}

// Permission returns Permission value that is concatenated
// using entity properties like subsystem, module and action.
func (pe *permissionEntity) Permission() charon.Permission {
	return charon.Permission(pe.Subsystem + ":" + pe.Module + ":" + pe.Action)
}

// PermissionRepository ...
type PermissionRepository interface {
	FindByUserID(int64) ([]*permissionEntity, error)
}

type permissionRepository struct {
	db *sql.DB
}

func newPermissionRepository(dbPool *sql.DB) *permissionRepository {
	return &permissionRepository{
		db: dbPool,
	}
}

// FindOneByID retrieves all permissions for user represented by given id.
func (pr *permissionRepository) FindByUserID(userID int64) ([]*permissionEntity, error) {
	query := `
		SELECT DISTINCT ON (p.id)
			` + tablePermissionColumns + `
		FROM ` + tablePermission + ` AS p
		LEFT JOIN charon.user_permissions AS up ON up.permission_id = p.id AND up.user_id = $1
		LEFT JOIN charon.user_groups AS ug ON ug.user_id = $1
		LEFT JOIN charon.group_permissions AS gp ON gp.permission_id = p.id AND gp.group_id = ug.group_id
		WHERE up.user_id = $1 OR ug.user_id = $1
	`

	rows, err := pr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	permissions := []*permissionEntity{}
	for rows.Next() {
		var p permissionEntity
		err = rows.Scan(
			&p.ID,
			&p.SubsystemID,
			&p.Subsystem,
			&p.Module,
			&p.Action,
			&p.CreatedAt,
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
