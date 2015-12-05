package main

import "github.com/piotrkowalczuk/pqcnstr"

const (
	tableGroup                                                 = "charon.group"
	tableGroupConstraintPrimaryKey          pqcnstr.Constraint = tableGroup + "_pkey"
	tableGroupConstraintUniqueName          pqcnstr.Constraint = tableGroup + "_name_key"
	tableGroupConstraintForeignKeyCreatedBy pqcnstr.Constraint = tableGroup + "_created_by_fkey"
	tableGroupConstraintForeignKeyUpdatedBy pqcnstr.Constraint = tableGroup + "_updated_by_fkey"
	tableGroupCreate                                           = `
		CREATE TABLE IF NOT EXISTS ` + tableGroup + ` (
			id         SERIAL,
			name       TEXT                      NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
			created_by INTEGER,
			updated_at TIMESTAMPTZ,
			updated_by INTEGER,

			CONSTRAINT "` + tableGroupConstraintPrimaryKey + `" PRIMARY KEY (id),
			CONSTRAINT "` + tableGroupConstraintUniqueName + `" UNIQUE (name),
			CONSTRAINT "` + tableUserConstraintForeignKeyCreatedBy + `" FOREIGN KEY (created_by) REFERENCES ` + tableUser + ` (id),
			CONSTRAINT "` + tableUserConstraintForeignKeyUpdatedBy + `" FOREIGN KEY (updated_by) REFERENCES ` + tableUser + ` (id)
		)
	`
	tableGroupColumns = `
		id, name, created_at, created_by, updated_at, updated_by
	`
)
