package main

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
)

type UserPermissionsRepository interface {
	Insert(entity *userPermissionsEntity) (*userPermissionsEntity, error)
	Exists(userID int64, permission charon.Permission) (bool, error)
}

func newUserPermissionsRepository(dbPool *sql.DB) UserPermissionsRepository {
	return &userPermissionsRepository{
		db:      dbPool,
		table:   tableUserPermissions,
		columns: tableUserPermissionsColumns,
	}
}

// Exists implements UserPermissionsRepository interface.
func (ugr *userPermissionsRepository) Exists(userID int64, permission charon.Permission) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM ` + tablePermission + ` AS p
			JOIN ` + ugr.table + ` AS up
				ON up.` + tableUserPermissionsColumnPermissionID + ` = p.` + tablePermissionColumnID + `
			WHERE p.` + tablePermissionColumnSubsystem + ` = $1
				AND p.` + tablePermissionColumnModule + ` = $2
				AND p.` + tablePermissionColumnAction + ` = $3
				AND up.` + tableUserPermissionsColumnUserID + ` = $4
		)
	`

	subsystem, module, action := permission.Split()
	var exists bool
	if err := ugr.db.QueryRow(
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
