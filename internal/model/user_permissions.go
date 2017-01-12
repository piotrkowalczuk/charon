package model

import (
	"context"
	"database/sql"
	"fmt"
)

// UserPermissionsProvider ...
type UserPermissionsProvider interface {
	Insert(context.Context, *UserPermissionsEntity) (*UserPermissionsEntity, error)
	DeleteByUserID(context.Context, int64) (int64, error)
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
			DB:      dbPool,
			Table:   TableUserPermissions,
			Columns: TableUserPermissionsColumns,
		},
		deleteByUserIDQuery: fmt.Sprintf("DELETE FROM %s WHERE %s = $1", TableUserPermissions, TableUserPermissionsColumnUserID),
	}
}

// DeleteByUserID removes all permissions of given user.
func (upr *UserPermissionsRepository) DeleteByUserID(ctx context.Context, id int64) (int64, error) {
	res, err := upr.DB.ExecContext(ctx, upr.deleteByUserIDQuery, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
