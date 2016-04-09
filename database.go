package charon

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/piotrkowalczuk/pqcomp"
)

type repositories struct {
	user             userProvider
	userGroups       userGroupsProvider
	userPermissions  userPermissionsProvider
	permission       permissionProvider
	group            groupProvider
	groupPermissions groupPermissionsProvider
}

func newRepositories(db *sql.DB) repositories {
	return repositories{
		user:             newUserRepository(db),
		userGroups:       newUserGroupsRepository(db),
		userPermissions:  newUserPermissionsRepository(db),
		permission:       newPermissionRepository(db),
		group:            newGroupRepository(db),
		groupPermissions: newGroupPermissionsRepository(db),
	}
}

func execQueries(db *sql.DB, queries ...string) (err error) {
	exec := func(query string) {
		if err != nil {
			return
		}

		_, err = db.Exec(query)
	}

	for _, q := range queries {
		exec(q)
	}

	return
}

func setupDatabase(db *sql.DB) error {
	return execQueries(
		db,
		schemaSQL,
	)
}

func teardownDatabase(db *sql.DB) error {
	return execQueries(
		db,
		`DROP SCHEMA IF EXISTS charon CASCADE`,
	)
}

func columns(names []string, prefix string) string {
	b := bytes.NewBuffer(nil)
	for i, n := range names {
		if i != 0 {
			b.WriteRune(',')
		}
		b.WriteString(prefix)
		b.WriteRune('.')
		b.WriteString(n)
	}

	return b.String()
}

func findQueryComp(db *sql.DB, table string, root, where *pqcomp.Composer, sort map[string]bool, columns []string) (*sql.Rows, error) {
	b := bytes.NewBufferString(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + table)

	if where.Len() != 0 {
		b.WriteString(` WHERE `)
		for where.Next() {
			if !where.First() {
				b.WriteString(" AND ")
			}

			fmt.Fprintf(b, "%s %s %s", where.Key(), where.Oper(), where.PlaceHolder())
		}
	}

	i := 0
SortLoop:
	for column, asc := range sort {
		if i != 0 {
			b.WriteString(", ")
		} else {
			b.WriteString(" ORDER BY ")
		}
		i++
		if asc {
			fmt.Fprintf(b, "%s ASC", column)
			continue SortLoop
		}

		fmt.Fprintf(b, "%s DESC ", column)
	}
	b.WriteString(" OFFSET $1 LIMIT $2")

	return db.Query(b.String(), root.Args()...)

}

func insertQueryComp(db *sql.DB, table string, insert *pqcomp.Composer, col []string) *sql.Row {
	b := bytes.NewBufferString(`INSERT INTO ` + table)

	if insert.Len() != 0 {
		b.WriteString(` (`)
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.Key())
		}
		insert.Reset()
		b.WriteString(`) VALUES (`)
		for insert.Next() {
			if !insert.First() {
				b.WriteString(", ")
			}

			fmt.Fprintf(b, "%s", insert.PlaceHolder())
		}
		b.WriteString(`)`)
		if len(col) > 0 {
			b.WriteString(" RETURNING ")
			b.WriteString(strings.Join(col, ","))
		}
	}

	return db.QueryRow(b.String(), insert.Args()...)
}

func existsManyToManyQuery(table, column1, column2 string) string {
	return `
		SELECT EXISTS(
			SELECT 1 FROM  ` + table + ` AS t
			WHERE t.` + column1 + ` = $1
				AND t.` + column2 + `= $2
		)
	`
}

func isGrantedQuery(table, columnID, columnSubsystem, columnModule, columnAction string) string {
	return `
		SELECT EXISTS(
			SELECT 1 FROM  ` + table + ` AS t
			WHERE t.` + columnID + ` = $1
				AND t.` + columnSubsystem + `= $2
				AND t.` + columnModule + `= $3
				AND t.` + columnAction + `= $4
		)
	`
}

func setManyToMany(db *sql.DB, table, column1, column2 string, id int64, ids []int64) (int64, int64, error) {
	var (
		err                    error
		aff, inserted, deleted int64
		tx                     *sql.Tx
		insert, exists         *sql.Stmt
		res                    sql.Result
		in                     []int64
		granted                bool
	)

	tx, err = db.Begin()
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

	if len(ids) > 0 {
		insert, err = tx.Prepare(`INSERT INTO ` + table + ` (` + column1 + `, ` + column2 + `) VALUES ($1, $2)`)
		if err != nil {
			return 0, 0, err
		}
		exists, err = tx.Prepare(existsManyToManyQuery(table, column1, column2))
		if err != nil {
			return 0, 0, err
		}

		in = make([]int64, 0, len(ids))
	InsertLoop:
		for _, idd := range ids {
			if err = exists.QueryRow(id, idd).Scan(&granted); err != nil {
				return 0, 0, err
			}
			// Given combination already exists, ignore.
			if granted {
				in = append(in, idd)
				granted = false
				continue InsertLoop
			}
			res, err = insert.Exec(id, idd)
			if err != nil {
				return 0, 0, err
			}

			aff, err = res.RowsAffected()
			if err != nil {
				return 0, 0, err
			}
			inserted += aff

			in = append(in, idd)
		}
	}

	delete := pqcomp.New(1, len(in))
	delete.AddArg(id)
	for _, id := range in {
		delete.AddExpr(column2, "IN", id)
	}

	query := bytes.NewBufferString(`DELETE FROM ` + table + ` WHERE ` + column1 + ` = $1`)
	for delete.Next() {
		if delete.First() {
			fmt.Fprint(query, " AND "+column2+" NOT IN (")
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

func setPermissions(db *sql.DB, table, columnID, columnSubsystem, columnModule, columnAction string, id int64, permissions Permissions) (int64, int64, error) {
	if len(permissions) == 0 {
		return 0, 0, errors.New("charon: permission cannot be set, none provided")
	}
	var (
		err                    error
		aff, inserted, deleted int64
		tx                     *sql.Tx
		insert, exists         *sql.Stmt
		res                    sql.Result
		in                     []Permission
		granted                bool
	)

	tx, err = db.Begin()
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

	var (
		subsystem, module, action string
	)

	if len(permissions) > 0 {
		insert, err = tx.Prepare(`INSERT INTO ` + table + ` (` + columnID + `, ` + columnSubsystem + `, ` + columnModule + `,` + columnAction + `) VALUES ($1, $2, $3, $4)`)
		if err != nil {
			return 0, 0, err
		}
		exists, err = tx.Prepare(isGrantedQuery(table, columnID, columnSubsystem, columnModule, columnAction))
		if err != nil {
			return 0, 0, err
		}

		in = make(Permissions, 0, len(permissions))
	InsertLoop:
		for _, p := range permissions {
			subsystem, module, action = p.Split()

			if err = exists.QueryRow(id, subsystem, module, action).Scan(&granted); err != nil {
				return 0, 0, fmt.Errorf("charon: error on permission check: %s", err.Error())
			}
			// Given combination already exists, ignore.
			if granted {
				in = append(in, p)
				granted = false
				continue InsertLoop
			}
			res, err = insert.Exec(id, subsystem, module, action)
			if err != nil {
				return 0, 0, fmt.Errorf("charon: error on permission insert: %s", err.Error())
			}

			aff, err = res.RowsAffected()
			if err != nil {
				return 0, 0, err
			}
			inserted += aff

			in = append(in, p)
		}
	}

	delete := pqcomp.New(1, len(in)*3)
	delete.AddArg(id)
	for _, p := range in {
		subsystem, module, action = p.Split()
		delete.AddExpr(columnSubsystem, "IN", subsystem)
		delete.AddExpr(columnModule, "IN", module)
		delete.AddExpr(columnAction, "IN", action)
	}

	query := bytes.NewBufferString(`DELETE FROM ` + table + ` WHERE ` + columnID + ` = $1`)
	if len(in) > 0 {
		fmt.Fprint(query, ` AND (`+columnSubsystem+`, `+columnModule+`, `+columnAction+`) NOT IN (`)
		for i, _ := range in {
			if i != 0 {
				fmt.Fprint(query, ", ")
			}
			delete.Next()
			fmt.Fprintf(query, "(%s", delete.PlaceHolder())
			delete.Next()
			fmt.Fprintf(query, ",%s,", delete.PlaceHolder())
			delete.Next()
			fmt.Fprintf(query, "%s)", delete.PlaceHolder())
		}
		fmt.Fprint(query, ")")
	}

	res, err = tx.Exec(query.String(), delete.Args()...)
	if err != nil {
		return 0, 0, fmt.Errorf("charon: error on redundant permission removal: %s", err.Error())
	}
	deleted, err = res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return inserted, deleted, nil
}
