package model

import "database/sql"

// GroupPermissionsProvider ...
type GroupPermissionsProvider interface {
	Insert(entity *GroupPermissionsEntity) (*GroupPermissionsEntity, error)
}

// GroupPermissionsRepository extends GroupPermissionsRepositoryBase
type GroupPermissionsRepository struct {
	GroupPermissionsRepositoryBase
}

// NewGroupPermissionsRepository ...
func NewGroupPermissionsRepository(dbPool *sql.DB) GroupPermissionsProvider {
	return &GroupPermissionsRepository{
		GroupPermissionsRepositoryBase{
			db:      dbPool,
			table:   TableGroupPermissions,
			columns: TableGroupPermissionsColumns,
		},
	}
}
