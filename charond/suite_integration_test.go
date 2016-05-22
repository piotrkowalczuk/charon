package charond

import (
	"io"
	"net"
	"testing"

	"io/ioutil"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charontest"
	"github.com/piotrkowalczuk/sklog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type integrationSuite struct {
	logger     log.Logger
	listener   net.Listener
	grpcServer *grpc.Server
	client     charon.RPCClient
	conn       *grpc.ClientConn
	server     charon.RPCServer

	mock struct {
		auth             *charontest.Charon
		user             *mockUserProvider
		group            *mockGroupProvider
		permission       *mockPermissionProvider
		userGroups       *mockUserGroupsProvider
		userPermissions  *mockUserPermissionsProvider
		groupPermissions *mockGroupPermissionsProvider
	}
}

func setupIntegrationSuite(t *testing.T) (*integrationSuite, error) {
	if testing.Short() {
		t.Skipf("integration suite ignored in short mode")
	}

	var (
		charonMock               charontest.Charon
		userRepoMock             mockUserProvider
		groupRepoMock            mockGroupProvider
		permissionRepoMock       mockPermissionProvider
		userGroupsRepoMock       mockUserGroupsProvider
		userPermissionsRepoMock  mockUserPermissionsProvider
		groupPermissionsRepoMock mockGroupPermissionsProvider
	)

	logger := log.NewJSONLogger(ioutil.Discard)
	rpcServer := &rpcServer{
		logger: logger,
		repository: repositories{
			user:             &userRepoMock,
			group:            &groupRepoMock,
			permission:       &permissionRepoMock,
			userGroups:       &userGroupsRepoMock,
			userPermissions:  &userPermissionsRepoMock,
			groupPermissions: &groupPermissionsRepoMock,
		},
	}
	is := integrationSuite{
		logger: logger,
		mock: struct {
			auth             *charontest.Charon
			user             *mockUserProvider
			group            *mockGroupProvider
			permission       *mockPermissionProvider
			userGroups       *mockUserGroupsProvider
			userPermissions  *mockUserPermissionsProvider
			groupPermissions *mockGroupPermissionsProvider
		}{
			auth:             &charonMock,
			user:             &userRepoMock,
			group:            &groupRepoMock,
			permission:       &permissionRepoMock,
			userGroups:       &userGroupsRepoMock,
			userPermissions:  &userPermissionsRepoMock,
			groupPermissions: &groupPermissionsRepoMock,
		},
		server: rpcServer,
	}

	err := is.setup(t, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &is, nil
}

func (is *integrationSuite) setup(t *testing.T, dialOpts ...grpc.DialOption) (err error) {
	is.listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return
	}

	grpclog.SetLogger(sklog.NewGRPCLogger(is.logger))
	var opts []grpc.ServerOption
	is.grpcServer = grpc.NewServer(opts...)

	charon.RegisterRPCServer(is.grpcServer, is.server)

	go is.grpcServer.Serve(is.listener)

	is.conn, err = grpc.Dial(is.listener.Addr().String(), dialOpts...)
	if err != nil {
		return err
	}

	is.client = charon.NewRPCClient(is.conn)

	return
}

func (is *integrationSuite) teardown(t *testing.T) (err error) {
	close := func(c io.Closer) {
		if err != nil {
			return
		}

		if c == nil {
			return
		}

		err = c.Close()
	}

	close(is.conn)
	is.grpcServer.Stop()
	is.listener.Close()

	return
}
