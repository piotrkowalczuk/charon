package main

import "database/sql"

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
		`CREATE SCHEMA IF NOT EXISTS charon`,
		string(tableUserCreate),
		string(tableGroupCreate),
		string(tablePermissionCreate),
		string(tableUserGroupsCreate),
		string(tableUserPermissionsCreate),
		string(tableGroupPermissionsCreate),
	)
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
