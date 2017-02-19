package charond

import (
	"flag"
	"os"
	"testing"

	"context"
	"time"

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

func getStringEnvOr(env, or string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return or
}

func timeout(ctx context.Context) context.Context {
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	return ctx
}
