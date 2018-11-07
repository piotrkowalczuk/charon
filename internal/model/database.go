package model

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
)

type repositories struct {
	user             UserProvider
	userGroups       UserGroupsProvider
	userPermissions  UserPermissionsProvider
	permission       PermissionProvider
	group            GroupProvider
	groupPermissions GroupPermissionsProvider
}

func newRepositories(db *sql.DB) repositories {
	return repositories{
		user:             NewUserRepository(db),
		userGroups:       NewUserGroupsRepository(db),
		userPermissions:  NewUserPermissionsRepository(db),
		permission:       NewPermissionRepository(db),
		group:            NewGroupRepository(db),
		groupPermissions: NewGroupPermissionsRepository(db),
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
		SQL,
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

func setManyToMany(db *sql.DB, ctx context.Context, table, column1, column2 string, id int64, ids []int64) (int64, int64, error) {
	var (
		err                    error
		aff, inserted, deleted int64
		tx                     *sql.Tx
		insert, exists         *sql.Stmt
		res                    sql.Result
		in                     []int64
		granted                bool
	)

	tx, err = db.BeginTx(ctx, nil)
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
		insert, err = tx.PrepareContext(ctx, `INSERT INTO `+table+` (`+column1+`, `+column2+`) VALUES ($1, $2)`)
		if err != nil {
			return 0, 0, err
		}
		exists, err = tx.PrepareContext(ctx, existsManyToManyQuery(table, column1, column2))
		if err != nil {
			return 0, 0, err
		}

		in = make([]int64, 0, len(ids))
	InsertLoop:
		for _, idd := range ids {
			if err = exists.QueryRowContext(ctx, id, idd).Scan(&granted); err != nil {
				return 0, 0, err
			}
			// Given combination already exists, ignore.
			if granted {
				in = append(in, idd)
				granted = false
				continue InsertLoop
			}
			res, err = insert.ExecContext(ctx, id, idd)
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

	delete := NewComposer(1)
	delete.WriteString("DELETE FROM ")
	delete.WriteString(table)
	delete.WriteString(" WHERE ")
	delete.WriteString(column1)
	delete.WriteString(" = ")
	if err = delete.WritePlaceholder(); err != nil {
		return 0, 0, err
	}
	delete.Add(id)
	if len(in) > 0 {
		delete.WriteString(" AND ")
		delete.WriteString(column2)
		delete.WriteString(" NOT IN (")
		for i, v := range in {
			if i != 0 {
				if _, err = delete.WriteString(","); err != nil {
					return 0, 0, err
				}
			}
			if err = delete.WritePlaceholder(); err != nil {
				return 0, 0, err
			}

			delete.Add(v)
		}
		if _, err = delete.WriteString(")"); err != nil {
			return 0, 0, err
		}
	}

	res, err = tx.ExecContext(ctx, delete.String(), delete.Args()...)
	if err != nil {
		return 0, 0, err
	}
	deleted, err = res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return inserted, deleted, nil
}

func setPermissions(db *sql.DB, ctx context.Context, table, columnID, columnSubsystem, columnModule, columnAction string, id int64, permissions charon.Permissions) (int64, int64, error) {
	if len(permissions) == 0 {
		return 0, 0, errors.New("permission cannot be set, none provided")
	}
	var (
		err                    error
		aff, inserted, deleted int64
		tx                     *sql.Tx
		insert, exists         *sql.Stmt
		res                    sql.Result
		in                     []charon.Permission
		granted                bool
	)

	tx, err = db.BeginTx(ctx, nil)
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

		in = make(charon.Permissions, 0, len(permissions))
	InsertLoop:
		for _, p := range permissions {
			subsystem, module, action = p.Split()

			if err = exists.QueryRow(id, subsystem, module, action).Scan(&granted); err != nil {
				return 0, 0, pqErrorPrefix(err, "error on permission check")
			}
			// Given combination already exists, ignore.
			if granted {
				in = append(in, p)
				granted = false
				continue InsertLoop
			}
			res, err = insert.Exec(id, subsystem, module, action)
			if err != nil {
				return 0, 0, pqErrorPrefix(err, "error on permission insert")
			}

			aff, err = res.RowsAffected()
			if err != nil {
				return 0, 0, err
			}
			inserted += aff

			in = append(in, p)
		}
	}

	delete := NewComposer(1)
	delete.WriteString("DELETE FROM ")
	delete.WriteString(table)
	delete.WriteString(" WHERE ")
	delete.WriteString(columnID)
	delete.WriteString(" = ")
	if err = delete.WritePlaceholder(); err != nil {
		return 0, 0, err
	}
	delete.Add(id)
	if len(in) > 0 {
		delete.WriteString(" AND (")
		delete.WriteString(columnSubsystem)
		delete.WriteString(",")
		delete.WriteString(columnModule)
		delete.WriteString(",")
		delete.WriteString(columnAction)
		delete.WriteString(") NOT IN (")
		for i, v := range in {
			if i != 0 {
				if _, err = delete.WriteString(","); err != nil {
					return 0, 0, err
				}
			}

			subsystem, module, action = v.Split()
			if _, err = delete.WriteString("("); err != nil {
				return 0, 0, err
			}
			if err = delete.WritePlaceholder(); err != nil {
				return 0, 0, err
			}
			delete.Add(subsystem)
			if _, err = delete.WriteString(","); err != nil {
				return 0, 0, err
			}

			if err = delete.WritePlaceholder(); err != nil {
				return 0, 0, err
			}
			delete.Add(module)
			if _, err = delete.WriteString(","); err != nil {
				return 0, 0, err
			}
			if err = delete.WritePlaceholder(); err != nil {
				return 0, 0, err
			}
			delete.Add(action)

			if _, err = delete.WriteString(")"); err != nil {
				return 0, 0, err
			}
		}
		if _, err = delete.WriteString(")"); err != nil {
			return 0, 0, err
		}
	}

	fmt.Println(delete.String())
	res, err = tx.Exec(delete.String(), delete.Args()...)
	if err != nil {
		return 0, 0, pqErrorPrefix(err, "error on redundant permission removal")
	}
	deleted, err = res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return inserted, deleted, nil
}

func pqErrorPrefix(err error, pre string) error {
	if pqerr, ok := err.(*pq.Error); ok {
		pqerr.Message = pre + ": " + pqerr.Message
		return pqerr
	}

	return errors.New(pre + ": " + err.Error())
}
