package model

import (
	"flag"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testPostgresAddress string
)

func TestMain(m *testing.M) {
	flag.StringVar(&testPostgresAddress, "postgres.address", getStringEnvOr("CHAROND_POSTGRES_ADDRESS", "postgres://localhost/test?sslmode=disable"), "")
	flag.Parse()

	os.Exit(m.Run())
}

type suite interface {
	setup(testing.T)
	teardown(testing.T)
}

func getStringEnvOr(env, or string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return or
}
