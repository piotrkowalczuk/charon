package charond

import (
	"database/sql"
	"io"
	"net"
	"os"
	"testing"
	"time"

	klog "github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type endToEndSuite struct {
	db     *sql.DB
	hasher charon.PasswordHasher

	charon       charon.RPCClient
	charonCloser io.Closer
	charonConn   *grpc.ClientConn

	mnemosyne       mnemosynerpc.SessionManagerClient
	mnemosyneConn   *grpc.ClientConn
	mnemosyneCloser io.Closer
}

func (etes *endToEndSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("e2e suite ignored in short mode")
	}

	var (
		err                       error
		mnemosyneAddr, charonAddr net.Addr
	)

	logger := sklog.NewTestLogger(t)
	_ = klog.NewJSONLogger(os.Stdout)
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))

	mnemosyneAddr, etes.mnemosyneCloser = mnemosyned.TestDaemon(t, mnemosyned.TestDaemonOpts{
		StoragePostgresAddress: testPostgresAddress,
	})
	t.Logf("mnemosyne deamon running on: %s", mnemosyneAddr.String())

	charonAddr, etes.charonCloser = TestDaemon(t, TestDaemonOpts{
		PostgresAddress:  testPostgresAddress,
		MnemosyneAddress: mnemosyneAddr.String(),
	})
	t.Logf("charon deamon running on: %s", charonAddr.String())

	if etes.mnemosyneConn, err = grpc.Dial(
		mnemosyneAddr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
	); err != nil {
		t.Fatalf("mnemosyne grpc connection error: %s", err.Error())
	}
	if etes.charonConn, err = grpc.Dial(
		charonAddr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
	); err != nil {
		t.Fatalf("charon grpc connection error: %s", err.Error())
	}
	if etes.db, err = sql.Open("postgres", testPostgresAddress); err != nil {
		t.Fatalf("postgres connection error: %s", err.Error())
	}
	if err := setupDatabase(etes.db); err != nil {
		t.Fatalf("database setup error: %s", err.Error())
	}
	if etes.hasher, err = charon.NewBCryptPasswordHasher(bcrypt.MinCost); err != nil {
		t.Fatalf("password hasher error: %s", err.Error())
	}

	etes.charon = charon.NewRPCClient(etes.charonConn)
	etes.mnemosyne = mnemosynerpc.NewSessionManagerClient(etes.mnemosyneConn)
}

func (etes *endToEndSuite) teardown(t *testing.T) {
	if err := teardownDatabase(etes.db); err != nil {
		t.Errorf("e2e suite database teardown error: %s", err.Error())
	}
	if err := etes.mnemosyneConn.Close(); err != nil {
		t.Errorf("e2e suite mnemosyne conn close error: %s", err.Error())
	}
	if err := etes.charonConn.Close(); err != nil {
		t.Errorf("e2e suite charon conn close error: %s", err.Error())
	}
	if err := etes.mnemosyneCloser.Close(); err != nil {
		t.Errorf("e2e suite mnemosyne closer close error: %s", err.Error())
	}
	if err := etes.charonCloser.Close(); err != nil {
		t.Errorf("e2e suite charon closer close error: %s", err.Error())
	}
	if err := etes.db.Close(); err != nil {
		t.Errorf("e2e suite database conn close error: %s", err.Error())
	}
}
