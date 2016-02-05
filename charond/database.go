package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/piotrkowalczuk/pqcomp"
)

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
	b, err := ioutil.ReadFile("../charond/schema.sql")
	if err != nil {
		b, err = ioutil.ReadFile("charond/schema.sql")
		if err != nil {
			return err
		}
	}

	_, err = db.Exec(string(b))
	return err
}

func tearDownDatabase(db *sql.DB) error {
	drop := func(tableName string) string {
		return "DROP TABLE IF EXISTS " + tableName + " CASCADE"
	}

	return execQueries(
		db,
		drop(tableUser),
		drop(tableGroup),
		drop(tablePermission),
		drop(tableUserGroups),
		drop(tableUserPermissions),
		drop(tableGroupPermissions),
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
			b.WriteString("RETURNING ")
			b.WriteString(strings.Join(col, ","))
		}
	}

	return db.QueryRow(b.String(), insert.Args()...)

}
