package main

import (
	"database/sql"
	"testing"
)

type postgresSuite struct {
	postgres   *sql.DB
	user       UserRepository
	permission PermissionRepository
	group      GroupRepository
}

func (ps *postgresSuite) Setup(t *testing.T) {
	config.parse()
	var err error

	ps.postgres, err = sql.Open("postgres", config.postgres.connectionString)
	if err != nil {
		t.Errorf("connection to postgres failed with error: %s", err.Error())
		t.FailNow()
	}

	setupDatabase(ps.postgres)

	ps.user = newUserRepository(ps.postgres)
}

func (ps *postgresSuite) Teardown(t *testing.T) {
	if err := tearDownDatabase(ps.postgres); err != nil {
		t.Errorf("unexpected error during database teardown: %s", err.Error())
	}

	ps.postgres.Close()
}
