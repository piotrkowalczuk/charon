package main

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
)

type UserPermissionsRepository interface {
	Insert(entity *userPermissionsEntity) (*userPermissionsEntity, error)
	IsGranted(userID int64, permission charon.Permission) (bool, error)
	Exists(userID, permissionID int64) (bool, error)
	Set(userID int64, permissionIDs []int64) (int64, int64, error)
}

func newUserPermissionsRepository(dbPool *sql.DB) UserPermissionsRepository {
	return &userPermissionsRepository{
		db:      dbPool,
		table:   tableUserPermissions,
		columns: tableUserPermissionsColumns,
	}
}

// IsGranted implements UserPermissionsRepository interface.
func (upr *userPermissionsRepository) IsGranted(userID int64, permission charon.Permission) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM ` + tablePermission + ` AS p
			JOIN ` + upr.table + ` AS up
				ON up.` + tableUserPermissionsColumnPermissionID + ` = p.` + tablePermissionColumnID + `
			WHERE p.` + tablePermissionColumnSubsystem + ` = $1
				AND p.` + tablePermissionColumnModule + ` = $2
				AND p.` + tablePermissionColumnAction + ` = $3
				AND up.` + tableUserPermissionsColumnUserID + ` = $4
		)
	`

	subsystem, module, action := permission.Split()
	var exists bool
	if err := upr.db.QueryRow(
		query,
		subsystem,
		module,
		action,
		userID,
	).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// Exists implements UserPermissionsRepository interface.
func (upr *userPermissionsRepository) Exists(userID, groupID int64) (bool, error) {
	var exists bool
	if err := upr.db.QueryRow(existsManyToManyQuery(
		upr.table,
		tableUserPermissionsColumnUserID,
		tableUserPermissionsColumnPermissionID,
	), userID, groupID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// Set implements UserPermissionsRepository interface.
func (upr *userPermissionsRepository) Set(userID int64, groupIDs []int64) (int64, int64, error) {
	return setManyToMany(upr.db, upr.table, tableUserPermissionsColumnUserID, tableUserPermissionsColumnPermissionID, userID, groupIDs)
}

