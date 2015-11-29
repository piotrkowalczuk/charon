package main

import "github.com/piotrkowalczuk/pqcnstr"

const (
	sqlCnstrPrimaryKeyGroup          pqcnstr.Constraint = "charon.group_pkey"
	sqlCnstrForeignKeyGroupCreatedBy pqcnstr.Constraint = "charon.group_created_by_fkey"
	sqlCnstrForeignKeyGroupUpdatedBy pqcnstr.Constraint = "charon.group_updated_by_fkey"
)
