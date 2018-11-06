package charond

import (
	"context"
	"flag"
	"os"
	"runtime/debug"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
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
	// TODO: leak
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

func assertError(t *testing.T, e1, e2 error) {
	if e1 != nil {
		if !grpcerr.Match(e1, e2) {
			t.Fatalf("error do not match, got %v", e2)
		}
	} else if e2 != nil {
		t.Fatal(e2)
	}
}

func recoverTest(t *testing.T) {
	t.Helper()

	if err := recover(); err != nil {
		t.Error(err, string(debug.Stack()))
	}
}

func brokenDate() time.Time {
	return time.Date(1, 1, 0, 0, 0, 0, 0, time.UTC)
}
