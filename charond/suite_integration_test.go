package charond

import (
	"io"
	"io/ioutil"
	"net"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/charontest"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/sklog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type integrationSuite struct {
	logger     log.Logger
	listener   net.Listener
	grpcServer *grpc.Server
	client     struct {
		auth       charonrpc.AuthClient
		user       charonrpc.UserManagerClient
		group      charonrpc.GroupManagerClient
		permission charonrpc.PermissionManagerClient
	}
	conn   *grpc.ClientConn
	server struct {
		auth       charonrpc.AuthServer
		user       charonrpc.UserManagerServer
		group      charonrpc.GroupManagerServer
		permission charonrpc.PermissionManagerServer
	}

	mock struct {
		auth             *charontest.Client
		user             *model.MockUserProvider
		group            *model.MockGroupProvider
		permission       *model.MockPermissionProvider
		userGroups       *model.MockUserGroupsProvider
		userPermissions  *model.MockUserPermissionsProvider
		groupPermissions *model.MockGroupPermissionsProvider
	}
}

func setupIntegrationSuite(t *testing.T) (*integrationSuite, error) {
	if testing.Short() {
		t.Skip("integration suite ignored in short mode")
	}

	var (
		charonMock               charontest.Client
		userRepoMock             model.MockUserProvider
		groupRepoMock            model.MockGroupProvider
		permissionRepoMock       model.MockPermissionProvider
		userGroupsRepoMock       model.MockUserGroupsProvider
		userPermissionsRepoMock  model.MockUserPermissionsProvider
		groupPermissionsRepoMock model.MockGroupPermissionsProvider
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
			auth             *charontest.Client
			user             *model.MockUserProvider
			group            *model.MockGroupProvider
			permission       *model.MockPermissionProvider
			userGroups       *model.MockUserGroupsProvider
			userPermissions  *model.MockUserPermissionsProvider
			groupPermissions *model.MockGroupPermissionsProvider
		}{
			auth:             &charonMock,
			user:             &userRepoMock,
			group:            &groupRepoMock,
			permission:       &permissionRepoMock,
			userGroups:       &userGroupsRepoMock,
			userPermissions:  &userPermissionsRepoMock,
			groupPermissions: &groupPermissionsRepoMock,
		},
		server: struct {
			auth       charonrpc.AuthServer
			user       charonrpc.UserManagerServer
			group      charonrpc.GroupManagerServer
			permission charonrpc.PermissionManagerServer
		}{
			auth:       newAuth(rpcServer),
			user:       newUserManager(rpcServer),
			group:      newGroupManager(rpcServer),
			permission: newPermissionManager(rpcServer),
		},
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

	charonrpc.RegisterAuthServer(is.grpcServer, is.server.auth)
	charonrpc.RegisterUserManagerServer(is.grpcServer, is.server.user)
	charonrpc.RegisterPermissionManagerServer(is.grpcServer, is.server.permission)
	charonrpc.RegisterGroupManagerServer(is.grpcServer, is.server.group)

	go is.grpcServer.Serve(is.listener)

	is.conn, err = grpc.Dial(is.listener.Addr().String(), dialOpts...)
	if err != nil {
		return err
	}

	is.client.auth = charonrpc.NewAuthClient(is.conn)
	is.client.user = charonrpc.NewUserManagerClient(is.conn)
	is.client.group = charonrpc.NewGroupManagerClient(is.conn)
	is.client.permission = charonrpc.NewPermissionManagerClient(is.conn)

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
