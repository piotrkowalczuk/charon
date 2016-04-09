package charon

import "database/sql"

type groupPermissionsProvider interface {
	Insert(entity *groupPermissionsEntity) (*groupPermissionsEntity, error)
}

func newGroupPermissionsRepository(dbPool *sql.DB) groupPermissionsProvider {
	return &groupPermissionsRepository{
		db:      dbPool,
		table:   tableGroupPermissions,
		columns: tableGroupPermissionsColumns,
	}
}
