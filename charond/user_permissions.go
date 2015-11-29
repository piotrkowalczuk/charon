package main

import "github.com/piotrkowalczuk/pqcnstr"

const (
	sqlCnstrForeignKeyUserPermissionsUserID       pqcnstr.Constraint = "charon.user_permissions_user_id_fkey"
	sqlCnstrForeignKeyUserPermissionsPermissionID pqcnstr.Constraint = "charon.user_permissions_permission_id_fkey"
	sqlCnstrForeignKeyUserPermissionsCreatedBy    pqcnstr.Constraint = "charon.user_permissions_created_by_fkey"
)
