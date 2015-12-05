package main

import "github.com/piotrkowalczuk/pqcnstr"

const (
	tableGroupPermissions                                                       = "charon.group_permissions"
	tableGroupPermissionsConstraintUniqueGroupIDPermissionID pqcnstr.Constraint = tableGroupPermissions + "_group_id_permission_id_key"
	tableGroupPermissionsConstraintForeignKeyGroupID         pqcnstr.Constraint = tableGroupPermissions + "_group_id_fkey"
	tableGroupPermissionsConstraintForeignKeyPermissionID    pqcnstr.Constraint = tableGroupPermissions + "_permission_id_fkey"
	tableGroupPermissionsConstraintForeignKeyCreatedBy       pqcnstr.Constraint = tableGroupPermissions + "_created_by_fkey"
	tableGroupPermissionsCreate                                                 = `
		CREATE TABLE IF NOT EXISTS ` + tableGroupPermissions + ` (
			group_id      INTEGER                   NOT NULL,
			permission_id INTEGER                   NOT NULL,
			created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
			created_by    INTEGER,

			CONSTRAINT "` + tableGroupPermissionsConstraintUniqueGroupIDPermissionID + `" UNIQUE (group_id, permission_id),
			CONSTRAINT "` + tableGroupPermissionsConstraintForeignKeyGroupID + `" FOREIGN KEY (group_id) REFERENCES ` + tableGroup + ` (id),
			CONSTRAINT "` + tableGroupPermissionsConstraintForeignKeyPermissionID + `" FOREIGN KEY (permission_id) REFERENCES ` + tablePermission + ` (id),
			CONSTRAINT "` + tableGroupPermissionsConstraintForeignKeyCreatedBy + `" FOREIGN KEY (created_by) REFERENCES ` + tableUser + ` (id)
		)
	`
	tableGroupPermissionsColumns = `
		group_id, permission_id, created_at, created_by
	`
)
