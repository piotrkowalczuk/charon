package main

import (
	"database/sql"

	"github.com/piotrkowalczuk/pqcomp"

	"bytes"
	"fmt"
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
	var exists bool
	if err := ugr.db.QueryRow(ugr.existsQuery(), userID, groupID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (ugr *userGroupsRepository) existsQuery() string {
	return `
		SELECT EXISTS(
			SELECT 1 FROM  ` + ugr.table + ` AS ug
			WHERE ug.` + tableUserGroupsColumnUserID + ` = $1
				AND ug.` + tableUserGroupsColumnGroupID + `= $2
		)
	`
}
func (ugr *userGroupsRepository) Set(userID int64, groupIDs []int64) (int64, int64, error) {
	var (
		err                    error
		aff, inserted, deleted int64
		tx                     *sql.Tx
		insert, exists         *sql.Stmt
		res                    sql.Result
		in                     []int64
		granted                bool
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

	if len(groupIDs) > 0 {
		insert, err = tx.Prepare(`INSERT INTO ` + ugr.table + ` (` + tableUserGroupsColumnUserID + `, ` + tableUserGroupsColumnGroupID + `) VALUES ($1, $2)`)
		if err != nil {
			return 0, 0, err
		}
		exists, err = tx.Prepare(ugr.existsQuery())
		if err != nil {
			return 0, 0, err
		}

		in = make([]int64, 0, len(groupIDs))
	InsertLoop:
		for _, groupID := range groupIDs {
			if err = exists.QueryRow(userID, groupID).Scan(&granted); err != nil {
				return 0, 0, err
			}
			// Given combination already exists, ignore.
			if granted {
				in = append(in, groupID)
				granted = false
				continue InsertLoop
			}
			res, err = insert.Exec(userID, groupID)
			if err != nil {
				return 0, 0, err
			}

			aff, err = res.RowsAffected()
			if err != nil {
				return 0, 0, err
			}
			inserted += aff

			in = append(in, groupID)
		}
	}

	delete := pqcomp.New(1, len(in))
	delete.AddArg(userID)
	for _, id := range in {
		delete.AddExpr(tableUserGroupsColumnGroupID, "IN", id)
	}

	query := bytes.NewBufferString(`DELETE FROM ` + ugr.table + ` WHERE ` + tableUserGroupsColumnUserID + ` = $1`)
	for delete.Next() {
		if delete.First() {
			fmt.Fprint(query, " AND "+tableUserGroupsColumnGroupID+" NOT IN (")
		} else {
			fmt.Fprint(query, ", ")

		}
		fmt.Fprintf(query, "%s", delete.PlaceHolder())
	}
	if len(in) > 0 {
		fmt.Fprint(query, ")")
	}

	res, err = tx.Exec(query.String(), delete.Args()...)
	if err != nil {
		return 0, 0, err
	}
	deleted, err = res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return inserted, deleted, nil
}
