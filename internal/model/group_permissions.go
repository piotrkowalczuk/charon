package model

import (
	"context"
	"database/sql"
)

// GroupPermissionsProvider ...
type GroupPermissionsProvider interface {
	Insert(context.Context, *GroupPermissionsEntity) (*GroupPermissionsEntity, error)
}

// GroupPermissionsRepository extends GroupPermissionsRepositoryBase
type GroupPermissionsRepository struct {
	GroupPermissionsRepositoryBase
}

// NewGroupPermissionsRepository ...
func NewGroupPermissionsRepository(dbPool *sql.DB) GroupPermissionsProvider {
	return &GroupPermissionsRepository{
		GroupPermissionsRepositoryBase{
			DB:      dbPool,
			Table:   TableGroupPermissions,
			Columns: TableGroupPermissionsColumns,
		},
	}
}
