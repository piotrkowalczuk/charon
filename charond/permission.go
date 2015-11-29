package main

import (
	"database/sql"
	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/pqcnstr"
)

const (
	sqlCnstrPrimaryKeyPermission                  pqcnstr.Constraint = "charon.permission_pkey"
	sqlCnstrUniquePermissionSubsystemModuleAction pqcnstr.Constraint = "charon.permission_subsystem_module_action_key"
	sqlCnstrForeignKeyPermissionSubsystemID       pqcnstr.Constraint = "charon.permission_subsystem_id_fkey"
	sqlCnstrForeignKeyPermissionCreatedBy         pqcnstr.Constraint = "charon.permission_created_by_fkey"
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

type userPermissionEntity struct {
	ID           int64
	UserID       int64
	PermissionID int64
	CreatedAt    *time.Time
	CreatedBy    int64
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
			p.id ,
			p.subsystem_id,
			p.subsystem,
			p.module,
			p.action,
			p.created_at
		FROM charon.permission AS p
		JOIN charon.user_permissions AS up ON up.user_id = $1
		JOIN charon.user_groups AS ug ON ug.user_id = $1
		JOIN charon.group_permissions AS gp ON gp.group_id = ug.group_id
	`

	rows, err := pr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	permissions := []*permissionEntity{}
	for rows.Next() {
		var p *permissionEntity
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

		permissions = append(permissions, p)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return permissions, nil
}
