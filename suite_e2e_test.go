package charon

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
	klog "github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type endToEndSuite struct {
	db             *sql.DB
	hasher         PasswordHasher
	userRepository userProvider

	charon       RPCClient
	charonConn   *grpc.ClientConn
	charonDaemon *Daemon

	mnemosyne       mnemosyne.RPCClient
	mnemosyneConn   *grpc.ClientConn
	mnemosyneDaemon *mnemosyne.Daemon
}

func (etes *endToEndSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("e2e suite ignored in short mode")
	}

	var err error

	mnemosyneTCP := listenTCP(t)
	charonTCP := listenTCP(t)
	logger := sklog.NewTestLogger(t)
	_ = klog.NewJSONLogger(os.Stdout)
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))

	etes.mnemosyneDaemon = mnemosyne.NewDaemon(&mnemosyne.DaemonOpts{
		Namespace:              "mnemosyne",
		MonitoringEngine:       mnemosyne.MonitoringEnginePrometheus,
		StoragePostgresAddress: testPostgresAddress,
		Logger:                 logger,
		RPCListener:            mnemosyneTCP,
	})
	if err = etes.mnemosyneDaemon.Run(); err != nil {
		t.Fatalf("mnemosyne daemon start error: %s", err.Error())
	}
	t.Logf("mnemosyne deamon running on: %s", etes.mnemosyneDaemon.Addr().String())

	etes.charonDaemon = NewDaemon(&DaemonOpts{
		Namespace:          "charon",
		MonitoringEngine:   MonitoringEnginePrometheus,
		MnemosyneAddress:   etes.mnemosyneDaemon.Addr().String(),
		Logger:             logger,
		PostgresAddress:    testPostgresAddress,
		RPCListener:        charonTCP,
		PasswordBCryptCost: bcrypt.MinCost,
	})
	if err = etes.charonDaemon.Run(); err != nil {
		t.Fatalf("charon daemon start error: %s", err.Error())
	}
	t.Logf("charon deamon running on: %s", etes.charonDaemon.Addr().String())

	if etes.mnemosyneConn, err = grpc.Dial(
		mnemosyneTCP.Addr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
	); err != nil {
		t.Fatalf("mnemosyne grpc connection error: %s", err.Error())
	}
	if etes.charonConn, err = grpc.Dial(
		charonTCP.Addr().String(),
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
	if etes.hasher, err = NewBCryptPasswordHasher(bcrypt.MinCost); err != nil {
		t.Fatalf("password hasher error: %s", err.Error())
	}

	etes.charon = NewRPCClient(etes.charonConn)
	etes.mnemosyne = mnemosyne.NewRPCClient(etes.mnemosyneConn)
	etes.userRepository = newUserRepository(etes.db)

	if _, err = createDumyTestUser(etes.userRepository, etes.hasher); err != nil {
		t.Fatalf("dummy user error: %s", err.Error())
	}
}

func (etes *endToEndSuite) teardown(t *testing.T) {
	grpcClose := func(conn *grpc.ClientConn) error {
		state, err := conn.State()
		if err != nil {
			return err
		}
		if state != grpc.Shutdown {
			if err = conn.Close(); err != nil {
				return err
			}
		}
		return nil
	}

	if err := teardownDatabase(etes.db); err != nil {
		t.Errorf("e2e suite database teardown error: %s", err.Error())
	}
	if err := grpcClose(etes.mnemosyneConn); err != nil {
		t.Errorf("e2e suite mnemosyne conn close error: %s", err.Error())
	}
	if err := grpcClose(etes.charonConn); err != nil {
		t.Errorf("e2e suite charon conn close error: %s", err.Error())
	}

	if err := etes.mnemosyneDaemon.Close(); err != nil {
		t.Errorf("e2e suite mnemosyne daemon close error: %s", err.Error())
	}
	if err := etes.charonDaemon.Close(); err != nil {
		t.Errorf("e2e suite charon daemon close error: %s", err.Error())
	}

	if err := etes.db.Close(); err != nil {
		t.Errorf("e2e suite database conn close error: %s", err.Error())
	}
}
