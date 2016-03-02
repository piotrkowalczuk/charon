package main

import "database/sql"

type UserPermissionsRepository interface {
	Insert(entity *userPermissionsEntity) (*userPermissionsEntity, error)
}

func newUserPermissionsRepository(dbPool *sql.DB) UserPermissionsRepository {
	return &userPermissionsRepository{
		db:      dbPool,
		table:   tableUserPermissions,
		columns: tableUserPermissionsColumns,
	}
}
