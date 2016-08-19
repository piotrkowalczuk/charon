package charond

import (
	"database/sql"
)

type userGroupsProvider interface {
	insert(entity *userGroupsEntity) (*userGroupsEntity, error)
	Exists(userID, groupID int64) (bool, error)
	find(criteria *userGroupsCriteria) ([]*userGroupsEntity, error)
	Set(userID int64, groupIDs []int64) (int64, int64, error)
}

type userGroupsRepository struct {
	userGroupsRepositoryBase
}

func newUserGroupsRepository(dbPool *sql.DB) userGroupsProvider {
	return &userGroupsRepository{
		userGroupsRepositoryBase: userGroupsRepositoryBase{
			db:      dbPool,
			table:   tableUserGroups,
			columns: tableUserGroupsColumns,
		},
	}
}

// Exists implements UserGroupsRepository interface.
func (ugr *userGroupsRepository) Exists(userID, groupID int64) (bool, error) {
	var exists bool
	if err := ugr.db.QueryRow(existsManyToManyQuery(ugr.table, tableUserGroupsColumnUserID, tableUserGroupsColumnGroupID), userID, groupID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// Set implements UserGroupsRepository interface.
func (ugr *userGroupsRepository) Set(userID int64, groupIDs []int64) (int64, int64, error) {
	return setManyToMany(ugr.db, ugr.table, tableUserGroupsColumnUserID, tableUserGroupsColumnGroupID, userID, groupIDs)
}
