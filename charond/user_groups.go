package main

import "github.com/piotrkowalczuk/pqcnstr"

const (
	sqlCnstrForeignKeyUserGroupsUserID    pqcnstr.Constraint = "charon.user_groups_user_id_fkey"
	sqlCnstrForeignKeyUserGroupsGroupID   pqcnstr.Constraint = "charon.user_groups_group_id_fkey"
	sqlCnstrForeignKeyUserGroupsCreatedBy pqcnstr.Constraint = "charon.user_groups_created_by_fkey"
)
