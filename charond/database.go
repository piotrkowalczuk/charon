package main

import (
	"database/sql"
	"io/ioutil"
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
