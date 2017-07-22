package charond

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func assertErrorCode(t *testing.T, err error, code codes.Code, msg string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != code {
			t.Fatalf("wrong error code, expected '%s' but got '%s' for error: %s", code, st.Code(), err.Error())
		}
		if st.Message() != msg {
			t.Fatalf("wrong error message, expected '%s' but got '%s' for error: %s", msg, st.Message(), err.Error())
		}
	} else {
		t.Fatalf("expected grpc error, got %T", err)
	}
}
