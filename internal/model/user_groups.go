package model

import (
	"database/sql"
)

// UserGroupsProvider ...
type UserGroupsProvider interface {
	Insert(entity *UserGroupsEntity) (*UserGroupsEntity, error)
	Exists(userID, groupID int64) (bool, error)
	Find(criteria *UserGroupsCriteria) ([]*UserGroupsEntity, error)
	Set(userID int64, groupIDs []int64) (int64, int64, error)
}

// UserGroupsRepository ...
type UserGroupsRepository struct {
	UserGroupsRepositoryBase
}

// NewUserGroupsRepository ...
func NewUserGroupsRepository(dbPool *sql.DB) UserGroupsProvider {
	return &UserGroupsRepository{
		UserGroupsRepositoryBase: UserGroupsRepositoryBase{
			db:      dbPool,
			table:   TableUserGroups,
			columns: TableUserGroupsColumns,
		},
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
