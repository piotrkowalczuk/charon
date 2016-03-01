package main

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
)

type GroupPermissionsRepository interface {
	Insert(entity *groupPermissionsEntity) (*groupPermissionsEntity, error)
	IsGranted(userID int64, permission charon.Permission) (bool, error)
	Exists(userID, permissionID int64) (bool, error)
	Set(userID int64, permissionIDs []int64) (int64, int64, error)
}

func newGroupPermissionsRepository(dbPool *sql.DB) GroupPermissionsRepository {
	return &groupPermissionsRepository{
		db:      dbPool,
		table:   tableGroupPermissions,
		columns: tableGroupPermissionsColumns,
	}
}

// IsGranted implements GroupPermissionsRepository interface.
func (upr *groupPermissionsRepository) IsGranted(groupID int64, permission charon.Permission) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM ` + tablePermission + ` AS p
			JOIN ` + upr.table + ` AS up
				ON up.` + tableGroupPermissionsColumnPermissionID + ` = p.` + tablePermissionColumnID + `
			WHERE p.` + tablePermissionColumnSubsystem + ` = $1
				AND p.` + tablePermissionColumnModule + ` = $2
				AND p.` + tablePermissionColumnAction + ` = $3
				AND up.` + tableGroupPermissionsColumnGroupID + ` = $4
		)
	`

	subsystem, module, action := permission.Split()
	var exists bool
	if err := upr.db.QueryRow(
		query,
		subsystem,
		module,
		action,
		groupID,
	).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// Exists implements GroupPermissionsRepository interface.
func (upr *groupPermissionsRepository) Exists(userID, groupID int64) (bool, error) {
	var exists bool
	if err := upr.db.QueryRow(existsManyToManyQuery(
		upr.table,
		tableGroupPermissionsColumnGroupID,
		tableGroupPermissionsColumnPermissionID,
	), userID, groupID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// Set implements GroupPermissionsRepository interface.
func (upr *groupPermissionsRepository) Set(userID int64, groupIDs []int64) (int64, int64, error) {
	return setManyToMany(upr.db, upr.table, tableGroupPermissionsColumnGroupID, tableGroupPermissionsColumnPermissionID, userID, groupIDs)
}

