package main

import (
	"database/sql"
	"strings"

	"strconv"
)

type UserGroupsRepository interface {
	Insert(entity *userGroupsEntity) (*userGroupsEntity, error)
	Exists(userID, groupID int64) (bool, error)
	Find(criteria *userGroupsCriteria) ([]*userGroupsEntity, error)
	Set(userID int64, groupIDs []int64) (int64, int64, error)
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

func (ugr *userGroupsRepository) Set(userID int64, groupIDs []int64) (int64, int64, error) {
	var (
		err                    error
		aff, inserted, deleted int64
		tx                     *sql.Tx
		stmt                   *sql.Stmt
		res                    sql.Result
		in                     []string
	)

	tx, err = ugr.db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt, err = tx.Prepare(`INSERT INTO ` + ugr.table + ` (` + tableUserGroupsColumnUserID + `, ` + tableUserGroupsColumnGroupID + `) VALUES ($1, $2)`)
	if err != nil {
		return 0, 0, err
	}

	in = make([]string, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		res, err = stmt.Exec(userID, groupID)
		if err != nil {
			return 0, 0, err
		}

		aff, err = res.RowsAffected()
		if err != nil {
			return 0, 0, err
		}
		inserted += aff

		in = append(in, strconv.FormatInt(groupID, 10))
	}

	res, err = tx.Exec(`DELETE FROM `+ugr.table+` WHERE `+tableUserGroupsColumnUserID+` == $1 AND `+tableUserGroupsColumnGroupID+` NOT IN ($2)`,
		userID,
		strings.Join(in, ","),
	)
	if err != nil {
		return 0, 0, err
	}
	deleted, err = res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return inserted, deleted, nil
}
