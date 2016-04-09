package charon

import "database/sql"

type userPermissionsProvider interface {
	Insert(entity *userPermissionsEntity) (*userPermissionsEntity, error)
}

func newUserPermissionsRepository(dbPool *sql.DB) userPermissionsProvider {
	return &userPermissionsRepository{
		db:      dbPool,
		table:   tableUserPermissions,
		columns: tableUserPermissionsColumns,
	}
}
