package main

import (
	"time"

	"github.com/piotrkowalczuk/pqcnstr"
)

const (
	tableUserGroups                                                 = "charon.user_groups"
	tableUserGroupsConstraintUniqueUserIDGroupID pqcnstr.Constraint = tableUserGroups + "_user_id_group_id_key"
	tableUserGroupsConstraintForeignKeyUserID    pqcnstr.Constraint = tableUserGroups + "_user_id_fkey"
	tableUserGroupsConstraintForeignKeyGroupID   pqcnstr.Constraint = tableUserGroups + "_group_id_fkey"
	tableUserGroupsConstraintForeignKeyCreatedBy pqcnstr.Constraint = tableUserGroups + "_created_by_fkey"
	tableUserGroupsCreate                                           = `
		CREATE TABLE IF NOT EXISTS ` + tableUserGroups + ` (
			user_id    INTEGER                   NOT NULL,
			group_id   INTEGER                   NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
			created_by INTEGER,

			CONSTRAINT "` + tableUserGroupsConstraintUniqueUserIDGroupID + `" UNIQUE (user_id, group_id),
			CONSTRAINT "` + tableUserGroupsConstraintForeignKeyUserID + `" FOREIGN KEY (user_id) REFERENCES ` + tableUser + ` (id),
			CONSTRAINT "` + tableUserGroupsConstraintForeignKeyGroupID + `" FOREIGN KEY (group_id) REFERENCES ` + tableGroup + ` (id),
			CONSTRAINT "` + tableUserGroupsConstraintForeignKeyCreatedBy + `" FOREIGN KEY (created_by) REFERENCES ` + tableUser + ` (id)
		)
	`
	tableUserGroupsColumns = `
		user_id, group_id, created_at, created_by
	`
)

type userGroupEntity struct {
	UserID    int64
	GroupID   int64
	CreatedAt *time.Time
	CreatedBy int64
}
