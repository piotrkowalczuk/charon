package main

import (
	"database/sql"
	"time"
)

// PermissionEntity ...
type PermissionEntity struct {
	ID          int64
	SubsystemID int64
	Subsystem   string
	Module      string
	Action      string
	CreatedAt   *time.Time
}

// UserPermissionEntity ...
type UserPermissionEntity struct {
	ID           int64
	UserID       int64
	PermissionID int64
	CreatedAt    *time.Time
	CreatedBy    int64
}

// PermissionRepository ...
type PermissionRepository interface {
	FindByUserID(int64) ([]*PermissionEntity, error)
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
func (pr *permissionRepository) FindByUserID(userID int64) ([]*PermissionEntity, error) {
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
		JOIN charon.user_group AS ug ON ug.user_id = $1
		JOIN charon.group_permissions AS gp ON gp.group_id = ug.group_id
	`

	rows, err := pr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	permissions := []*PermissionEntity{}
	for rows.Next() {
		var p *PermissionEntity
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
