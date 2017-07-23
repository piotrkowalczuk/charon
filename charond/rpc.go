package charond

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type rpcServer struct {
	opts               DaemonOpts
	meta               metadata.MD
	logger             log.Logger
	ldap               *sync.Pool
	session            mnemosynerpc.SessionManagerClient
	passwordHasher     password.Hasher
	permissionRegistry model.PermissionRegistry
	repository         repositories
}

type auth struct {
	*actorHandler
	*loginHandler
	*logoutHandler
	*isGrantedHandler
	*isAuthenticatedHandler
	*belongsToHandler
}

func newAuth(server *rpcServer) *auth {
	return &auth{
		actorHandler:           &actorHandler{handler: newHandler(server)},
		belongsToHandler:       &belongsToHandler{handler: newHandler(server)},
		isGrantedHandler:       &isGrantedHandler{handler: newHandler(server)},
		isAuthenticatedHandler: &isAuthenticatedHandler{handler: newHandler(server)},
		loginHandler:           &loginHandler{handler: newHandler(server), hasher: server.passwordHasher, mappings: server.opts.LDAPMappings},
		logoutHandler:          &logoutHandler{handler: newHandler(server)},
	}
}

type userManager struct {
	*createUserHandler
	*deleteUserHandler
	*getUserHandler
	*listUserGroupsHandler
	*listUserPermissionsHandler
	*listUsersHandler
	*modifyUserHandler
	*setUserGroupsHandler
	*setUserPermissionsHandler
}

func newUserManager(server *rpcServer) *userManager {
	return &userManager{
		createUserHandler:          &createUserHandler{handler: newHandler(server), hasher: server.passwordHasher},
		deleteUserHandler:          &deleteUserHandler{handler: newHandler(server)},
		getUserHandler:             &getUserHandler{handler: newHandler(server)},
		listUserGroupsHandler:      &listUserGroupsHandler{handler: newHandler(server)},
		listUserPermissionsHandler: &listUserPermissionsHandler{handler: newHandler(server)},
		listUsersHandler:           &listUsersHandler{handler: newHandler(server)},
		modifyUserHandler:          &modifyUserHandler{handler: newHandler(server)},
		setUserGroupsHandler:       &setUserGroupsHandler{handler: newHandler(server)},
		setUserPermissionsHandler:  &setUserPermissionsHandler{handler: newHandler(server)},
	}
}

type permissionManager struct {
	*registerPermissionsHandler
	*getPermissionHandler
	*listPermissionsHandler
}

func newPermissionManager(server *rpcServer) *permissionManager {
	return &permissionManager{
		registerPermissionsHandler: &registerPermissionsHandler{handler: newHandler(server), registry: server.permissionRegistry},
		listPermissionsHandler:     &listPermissionsHandler{handler: newHandler(server)},
		getPermissionHandler:       &getPermissionHandler{handler: newHandler(server)},
	}
}

type groupManager struct {
	*getGroupHandler
	*deleteGroupHandler
	*modifyGroupHandler
	*listGroupsHandler
	*setGroupPermissionsHandler
	*createGroupHandler
	*listGroupPermissionsHandler
}

func newGroupManager(server *rpcServer) *groupManager {
	return &groupManager{
		getGroupHandler:             &getGroupHandler{handler: newHandler(server)},
		deleteGroupHandler:          &deleteGroupHandler{handler: newHandler(server)},
		modifyGroupHandler:          &modifyGroupHandler{handler: newHandler(server)},
		listGroupsHandler:           &listGroupsHandler{handler: newHandler(server)},
		setGroupPermissionsHandler:  &setGroupPermissionsHandler{handler: newHandler(server)},
		createGroupHandler:          &createGroupHandler{handler: newHandler(server)},
		listGroupPermissionsHandler: &listGroupPermissionsHandler{handler: newHandler(server)},
	}
}

func unaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		wrap := func(current grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return current(currentCtx, currentReq, info, next)
			}
		}
		chain := handler
		for _, i := range interceptors {
			chain = wrap(i, chain)
		}
		return chain(ctx, req)
	}
}
