package main

import "database/sql"

type GroupPermissionsRepository interface {
	Insert(entity *groupPermissionsEntity) (*groupPermissionsEntity, error)
}

func newGroupPermissionsRepository(dbPool *sql.DB) GroupPermissionsRepository {
	return &groupPermissionsRepository{
		db:      dbPool,
		table:   tableGroupPermissions,
		columns: tableGroupPermissionsColumns,
	}
}
