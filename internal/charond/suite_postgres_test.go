package charond

import (
	"database/sql"
	"testing"

	"go.uber.org/zap"
)

type postgresSuite struct {
	logger     *zap.Logger
	db         *sql.DB
	repository repositories
}

func (ps *postgresSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres suite ignored in short mode")
	}

	var err error

	ps.logger = zap.L()
	ps.db, err = initPostgres(testPostgresAddress, true, ps.logger)
	if err != nil {
		t.Fatalf("postgres connection (%s) error: %s", testPostgresAddress, err.Error())
	}

	ps.repository = newRepositories(ps.db)
}

func (ps *postgresSuite) teardown(t *testing.T) {
	var err error

	if err = teardownDatabase(ps.db); err != nil {
		t.Fatalf("postgres suite database teardown error: %s", err.Error())
	}
	if err = ps.db.Close(); err != nil {
		t.Fatalf("postgres suite teardown database connection error: %s", err.Error())
	}
}
