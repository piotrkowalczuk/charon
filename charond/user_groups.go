package main

import "database/sql"

type UserGroupsRepository interface {
	Insert(entity *userGroupsEntity) (*userGroupsEntity, error)
	Exists(userID, groupID int64) (bool, error)
}

func newUserGroupsRepository(dbPool *sql.DB) UserGroupsRepository {
	return &userGroupsRepository{
		db:      dbPool,
		table:   tableUserGroups,
		columns: tableUserGroupsColumns,
	}
}

// Exists implements UserGroupsRepository interface.
func (ugr *userGroupsRepository) Exists(userID, groupID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM  ` + tableUserGroups + ` AS ug
			WHERE ug.user_id = $1
				AND ug.group_id = $2
		)
	`
	var exists bool
	if err := ugr.db.QueryRow(query, userID, groupID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
