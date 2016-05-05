package charond

import (
	"database/sql"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
)

type postgresSuite struct {
	logger     log.Logger
	db         *sql.DB
	repository repositories
}

func (ps *postgresSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres suite ignored in short mode")
	}

	var err error

	ps.logger = sklog.NewTestLogger(t)
	ps.db, err = initPostgres(testPostgresAddress, EnvironmentTest, ps.logger)
	if err != nil {
		t.Fatalf("postgres connection (%s) error: %s", testPostgresAddress, err.Error())
	}

	ps.repository = newRepositories(ps.db)
}

func (ps *postgresSuite) teardown(t *testing.T) {
	var err error

	if err = ps.db.Close(); err != nil {
		t.Fatal(err)
	}
}

func assertf(t *testing.T, is bool, msg string, args ...interface{}) bool {
	if !is {
		t.Errorf(msg, args...)
	}

	return is
}

func assert(t *testing.T, is bool, msg string) bool {
	if !is {
		t.Errorf(msg)
	}

	return is
}

func assertfTime(t *testing.T, tm *time.Time, msg string, args ...interface{}) bool {
	if tm == nil || tm.IsZero() {
		t.Errorf(msg, args...)
		return false
	}

	return true
}
