package charond

import "database/sql"

type userPermissionsProvider interface {
	insert(entity *userPermissionsEntity) (*userPermissionsEntity, error)
}

type userPermissionsRepository struct {
	userPermissionsRepositoryBase
}

func newUserPermissionsRepository(dbPool *sql.DB) userPermissionsProvider {
	return &userPermissionsRepository{
		userPermissionsRepositoryBase: userPermissionsRepositoryBase{
			db:      dbPool,
			table:   tableUserPermissions,
			columns: tableUserPermissionsColumns,
		},
	}
}
