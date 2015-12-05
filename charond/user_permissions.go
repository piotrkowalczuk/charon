package main

import "github.com/piotrkowalczuk/pqcnstr"

const (
	tableUserPermissions                                                      = "charon.user_permissions"
	tableUserPermissionsConstraintUniqueUserIDPermissionID pqcnstr.Constraint = tableUserPermissions + "_user_id_permission_id_key"
	tableUserPermissionsConstraintForeignKeyUserID         pqcnstr.Constraint = tableUserPermissions + "_user_id_fkey"
	tableUserPermissionsConstraintForeignKeyPermissionID   pqcnstr.Constraint = tableUserPermissions + "_permission_id_fkey"
	tableUserPermissionsConstraintForeignKeyCreatedBy      pqcnstr.Constraint = tableUserPermissions + "_created_by_fkey"
	tableUserPermissionsCreate                                                = `
		CREATE TABLE IF NOT EXISTS ` + tableUserPermissions + ` (
			user_id       INTEGER                   NOT NULL,
			permission_id INTEGER                   NOT NULL,
			created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
			created_by    INTEGER,

			CONSTRAINT "` + tableUserPermissionsConstraintUniqueUserIDPermissionID + `" UNIQUE (user_id, permission_id),
			CONSTRAINT "` + tableUserPermissionsConstraintForeignKeyUserID + `" FOREIGN KEY (user_id) REFERENCES ` + tableUser + ` (id),
			CONSTRAINT "` + tableUserPermissionsConstraintForeignKeyPermissionID + `" FOREIGN KEY (permission_id) REFERENCES ` + tablePermission + ` (id),
			CONSTRAINT "` + tableUserPermissionsConstraintForeignKeyCreatedBy + `" FOREIGN KEY (created_by) REFERENCES ` + tableUser + ` (id)
		)
	`
	tableUserPermissionsColumns = `
		user_id, permission_id, created_at, created_by
	`
)
