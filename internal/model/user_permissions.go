package model

import "database/sql"

// UserPermissionsProvider ...
type UserPermissionsProvider interface {
	Insert(entity *UserPermissionsEntity) (*UserPermissionsEntity, error)
}

// UserPermissionsRepository extends UserPermissionsRepositoryBase
type UserPermissionsRepository struct {
	UserPermissionsRepositoryBase
}

// NewUserPermissionsRepository ...
func NewUserPermissionsRepository(dbPool *sql.DB) UserPermissionsProvider {
	return &UserPermissionsRepository{
		UserPermissionsRepositoryBase: UserPermissionsRepositoryBase{
			db:      dbPool,
			table:   TableUserPermissions,
			columns: TableUserPermissionsColumns,
		},
	}
}
