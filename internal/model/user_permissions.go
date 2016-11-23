package model

import (
	"database/sql"
	"fmt"
)

// UserPermissionsProvider ...
type UserPermissionsProvider interface {
	Insert(entity *UserPermissionsEntity) (*UserPermissionsEntity, error)
	DeleteByUserID(id int64) (int64, error)
}

// UserPermissionsRepository extends UserPermissionsRepositoryBase
type UserPermissionsRepository struct {
	UserPermissionsRepositoryBase
	deleteByUserIDQuery string
}

// NewUserPermissionsRepository ...
func NewUserPermissionsRepository(dbPool *sql.DB) UserPermissionsProvider {
	return &UserPermissionsRepository{
		UserPermissionsRepositoryBase: UserPermissionsRepositoryBase{
			db:      dbPool,
			table:   TableUserPermissions,
			columns: TableUserPermissionsColumns,
		},
		deleteByUserIDQuery: fmt.Sprintf("DELETE FROM %s WHERE %s = $1", TableUserPermissions, TableUserPermissionsColumnUserID),
	}
}

// DeleteByUserID removes all permissions of given user.
func (upr *UserPermissionsRepository) DeleteByUserID(id int64) (int64, error) {
	res, err := upr.db.Exec(upr.deleteByUserIDQuery, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
