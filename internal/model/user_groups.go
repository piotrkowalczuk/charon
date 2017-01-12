package model

import (
	"context"
	"database/sql"
	"fmt"
)

// UserGroupsProvider ...
type UserGroupsProvider interface {
	Insert(context.Context, *UserGroupsEntity) (*UserGroupsEntity, error)
	Exists(ctx context.Context, userID, groupID int64) (bool, error)
	Find(context.Context, *UserGroupsCriteria) ([]*UserGroupsEntity, error)
	Set(ctx context.Context, userID int64, groupIDs []int64) (int64, int64, error)
	DeleteByUserID(context.Context, int64) (int64, error)
}

// UserGroupsRepository ...
type UserGroupsRepository struct {
	UserGroupsRepositoryBase
	deleteByUserIDQuery string
}

// NewUserGroupsRepository ...
func NewUserGroupsRepository(dbPool *sql.DB) UserGroupsProvider {
	return &UserGroupsRepository{
		UserGroupsRepositoryBase: UserGroupsRepositoryBase{
			DB:      dbPool,
			Table:   TableUserGroups,
			Columns: TableUserGroupsColumns,
		},
		deleteByUserIDQuery: fmt.Sprintf("DELETE FROM %s WHERE %s = $1", TableUserGroups, TableUserGroupsColumnUserID),
	}
}

// Exists implements UserGroupsProvider interface.
func (ugr *UserGroupsRepository) Exists(ctx context.Context, userID, groupID int64) (bool, error) {
	var exists bool
	if err := ugr.DB.QueryRowContext(ctx, existsManyToManyQuery(ugr.Table, TableUserGroupsColumnUserID, TableUserGroupsColumnGroupID), userID, groupID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// Set implements UserGroupsProvider interface.
func (ugr *UserGroupsRepository) Set(ctx context.Context, userID int64, groupIDs []int64) (int64, int64, error) {
	return setManyToMany(ugr.DB, ctx, ugr.Table, TableUserGroupsColumnUserID, TableUserGroupsColumnGroupID, userID, groupIDs)
}

// DeleteByUserID removes user from all groups he belongs to.
func (ugr *UserGroupsRepository) DeleteByUserID(ctx context.Context, id int64) (int64, error) {
	res, err := ugr.DB.ExecContext(ctx, ugr.deleteByUserIDQuery, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
