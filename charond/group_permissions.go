package main

import "github.com/piotrkowalczuk/pqcnstr"

const (
	sqlCnstrForeignKeyGroupPermissionsUserID    pqcnstr.Constraint = "charon.user_permissions_user_id_fkey"
	sqlCnstrForeignKeyGroupPermissionsGroupID   pqcnstr.Constraint = "charon.user_permissions_group_id_fkey"
	sqlCnstrForeignKeyGroupPermissionsCreatedBy pqcnstr.Constraint = "charon.user_permissions_created_by_fkey"
)
