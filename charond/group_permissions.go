package charond

import "database/sql"

type groupPermissionsProvider interface {
	Insert(entity *groupPermissionsEntity) (*groupPermissionsEntity, error)
}

type groupPermissionsRepository struct {
	groupPermissionsRepositoryBase
}

func newGroupPermissionsRepository(dbPool *sql.DB) groupPermissionsProvider {
	return &groupPermissionsRepository{
		groupPermissionsRepositoryBase{
			db:      dbPool,
			table:   tableGroupPermissions,
			columns: tableGroupPermissionsColumns,
		},
	}
}
