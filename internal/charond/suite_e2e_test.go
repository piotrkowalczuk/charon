package charond

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/ntypes"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type endToEndSuite struct {
	db        *sql.DB
	hasher    password.Hasher
	userAgent string

	charon struct {
		auth         charonrpc.AuthClient
		user         charonrpc.UserManagerClient
		group        charonrpc.GroupManagerClient
		permission   charonrpc.PermissionManagerClient
		refreshToken charonrpc.RefreshTokenManagerClient
	}
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

	mnemosyneAddr, etes.mnemosyneCloser = mnemosyned.TestDaemon(t, mnemosyned.TestDaemonOpts{
		StoragePostgresAddress: testPostgresAddress,
	})
	t.Logf("mnemosyne deamon running on: %s", mnemosyneAddr.String())

	charonAddr, etes.charonCloser = TestDaemon(t, TestDaemonOpts{
		PostgresDebug:    true,
		PostgresAddress:  testPostgresAddress,
		MnemosyneAddress: mnemosyneAddr.String(),
	})
	t.Logf("charon deamon running on: %s", charonAddr.String())

	if etes.mnemosyneConn, err = grpc.Dial(
		mnemosyneAddr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
		grpc.WithUserAgent(etes.userAgent),
	); err != nil {
		t.Fatalf("mnemosyne grpc connection error: %s", err.Error())
	}
	if etes.charonConn, err = grpc.Dial(
		charonAddr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
		grpc.WithUserAgent(etes.userAgent),
	); err != nil {
		t.Fatalf("charon grpc connection error: %s", err.Error())
	}
	if etes.db, err = sql.Open("postgres", testPostgresAddress); err != nil {
		t.Fatalf("postgres connection error: %s", err.Error())
	}
	if err := setupDatabase(etes.db); err != nil {
		t.Fatalf("database setup error: %s", err.Error())
	}
	if etes.hasher, err = password.NewBCryptHasher(bcrypt.MinCost); err != nil {
		t.Fatalf("password hasher error: %s", err.Error())
	}

	etes.charon = struct {
		auth         charonrpc.AuthClient
		user         charonrpc.UserManagerClient
		group        charonrpc.GroupManagerClient
		permission   charonrpc.PermissionManagerClient
		refreshToken charonrpc.RefreshTokenManagerClient
	}{
		auth:         charonrpc.NewAuthClient(etes.charonConn),
		user:         charonrpc.NewUserManagerClient(etes.charonConn),
		group:        charonrpc.NewGroupManagerClient(etes.charonConn),
		permission:   charonrpc.NewPermissionManagerClient(etes.charonConn),
		refreshToken: charonrpc.NewRefreshTokenManagerClient(etes.charonConn),
	}
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

func (etes *endToEndSuite) createGroups(t *testing.T, ctx context.Context) ([]*charonrpc.Group, []int64) {
	var (
		ids    []int64
		groups []*charonrpc.Group
	)
	for i := 0; i < 10; i++ {
		res, err := etes.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
			Name: fmt.Sprintf("name-%d", i),
			Description: &ntypes.String{
				Valid: true,
				Chars: fmt.Sprintf("description-%d", i),
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		ids = append(ids, res.Group.Id)
		groups = append(groups, res.Group)
	}
	return groups, ids
}

func (etes *endToEndSuite) createUsers(t *testing.T, ctx context.Context) ([]*charonrpc.User, []int64) {
	var (
		ids   []int64
		users []*charonrpc.User
	)
	for i := 0; i < 10; i++ {
		res, err := etes.charon.user.Create(ctx, &charonrpc.CreateUserRequest{
			Username:      fmt.Sprintf("username-%d@example.com", i),
			FirstName:     fmt.Sprintf("first-name-%d", i),
			LastName:      fmt.Sprintf("last-name-%d", i),
			PlainPassword: fmt.Sprintf("password-%d", i),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		ids = append(ids, res.User.Id)
		users = append(users, res.User)
	}
	return users, ids
}
