package model

import (
	"database/sql"
	"fmt"
)

// UserGroupsProvider ...
type UserGroupsProvider interface {
	Insert(entity *UserGroupsEntity) (*UserGroupsEntity, error)
	Exists(userID, groupID int64) (bool, error)
	Find(criteria *UserGroupsCriteria) ([]*UserGroupsEntity, error)
	Set(userID int64, groupIDs []int64) (int64, int64, error)
	DeleteByUserID(id int64) (int64, error)
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
			db:      dbPool,
			table:   TableUserGroups,
			columns: TableUserGroupsColumns,
		},
		deleteByUserIDQuery: fmt.Sprintf("DELETE FROM %s WHERE %s = $1", TableUserGroups, TableUserGroupsColumnUserID),
	}
}

// Exists implements UserGroupsProvider interface.
func (ugr *UserGroupsRepository) Exists(userID, groupID int64) (bool, error) {
	var exists bool
	if err := ugr.db.QueryRow(existsManyToManyQuery(ugr.table, TableUserGroupsColumnUserID, TableUserGroupsColumnGroupID), userID, groupID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// Set implements UserGroupsProvider interface.
func (ugr *UserGroupsRepository) Set(userID int64, groupIDs []int64) (int64, int64, error) {
	return setManyToMany(ugr.db, ugr.table, TableUserGroupsColumnUserID, TableUserGroupsColumnGroupID, userID, groupIDs)
}

// DeleteByUserID removes user from all groups he belongs to.
func (ugr *UserGroupsRepository) DeleteByUserID(id int64) (int64, error) {
	res, err := ugr.db.Exec(ugr.deleteByUserIDQuery, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
