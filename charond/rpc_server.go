package main

import (
	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/pqcnstr"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type rpcServer struct {
	meta               metadata.MD
	logger             log.Logger
	monitor            *monitoring
	session            mnemosyne.Mnemosyne
	passwordHasher     charon.PasswordHasher
	permissionRegistry PermissionRegistry
	repository         struct {
		user       UserRepository
		permission PermissionRepository
		group      GroupRepository
	}
}

type actor struct {
	user        *userEntity
	session     *mnemosyne.Session
	permissions charon.Permissions
}

func (a *actor) isLocalhost() bool {
	return a.user == nil && a.session == nil && a.permissions == nil
}

// TODO: refactor
func mapUserError(err error) error {
	switch pqcnstr.FromError(err) {
	case tableUserConstraintPrimaryKey:
		return grpc.Errorf(codes.AlreadyExists, charon.ErrDescUserWithIDExists)
	case tableUserConstraintUsernameUnique:
		return grpc.Errorf(codes.AlreadyExists, charon.ErrDescUserWithUsernameExists)
	default:
		return err
	}
}

func (rs *rpcServer) loggerBackground(ctx context.Context, keyval ...interface{}) log.Logger {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return rs.logger
	}

	rid, ok := md["request_id"]
	if !ok {
		return rs.logger
	}

	return log.NewContext(rs.logger).With(append(keyval, "request_id", rid)...)
}

// Context create new context based on given metadata and instance metadata.
func (rs *rpcServer) Context(md metadata.MD) context.Context {
	if md.Len() == 0 {
		md = rs.meta
	} else {
		md = rs.metadata(md)
	}

	return metadata.NewContext(context.Background(), md)
}

func (rs *rpcServer) metadata(md metadata.MD) metadata.MD {
	for key, value := range rs.meta {
		if _, ok := md[key]; !ok {
			md[key] = value
		}
	}

	return md
}

// Login implements charon.RPCServer interface.
func (rs *rpcServer) Login(ctx context.Context, req *charon.LoginRequest) (*charon.LoginResponse, error) {
	h := &loginHandler{
		handler: newHandler(rs, ctx, "login"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "subject has been logged in")

	return resp, err
}

// Logout implements charon.RPCServer interface.
func (rs *rpcServer) Logout(ctx context.Context, req *charon.LogoutRequest) (*charon.LogoutResponse, error) {
	h := &logoutHandler{
		handler: newHandler(rs, ctx, "logout"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "subject has been logged out")

	return resp, err
}

// IsAuthenticated implements charon.RPCServer interface.
func (rs *rpcServer) IsAuthenticated(ctx context.Context, req *charon.IsAuthenticatedRequest) (*charon.IsAuthenticatedResponse, error) {
	h := &isAuthenticatedHandler{
		handler: newHandler(rs, ctx, "is_authenticated"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "subject has been logged out")

	return resp, err
}

// Subject implements charon.RPCServer interface.
func (rs *rpcServer) Subject(ctx context.Context, req *charon.SubjectRequest) (*charon.SubjectResponse, error) {
	h := &subjectHandler{
		handler: newHandler(rs, ctx, "subject"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "subject has been retrieved")

	return resp, err
}

// IsGranted implements charon.RPCServer interface.
func (rs *rpcServer) IsGranted(ctx context.Context, req *charon.IsGrantedRequest) (*charon.IsGrantedResponse, error) {
	h := &isGrantedHandler{
		handler: newHandler(rs, ctx, "is_granted"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "permission has been checked")

	return resp, err
}

// BelongsTo implements charon.RPCServer interface.
func (rs *rpcServer) BelongsTo(ctx context.Context, req *charon.BelongsToRequest) (*charon.BelongsToResponse, error) {
	h := &belongsToHandler{
		handler: newHandler(rs, ctx, "belongs_to"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "belonging to the group has been checked")

	return resp, err
}

// CreateGroup implements charon.RPCServer interface.
func (rs *rpcServer) CreateGroup(ctx context.Context, req *charon.CreateGroupRequest) (*charon.CreateGroupResponse, error) {
	h := &createGroupHandler{
		handler: newHandler(rs, ctx, "create_group"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "group has been created")

	return resp, err
}

// ModifyGroup implements charon.RPCServer interface.
func (rs *rpcServer) ModifyGroup(ctx context.Context, req *charon.ModifyGroupRequest) (*charon.ModifyGroupResponse, error) {
	h := &modifyGroupHandler{
		handler: newHandler(rs, ctx, "modify_group"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "group has been created")

	return resp, err
}

// DeleteGroup implements charon.RPCServer interface.
func (rs *rpcServer) DeleteGroup(ctx context.Context, req *charon.DeleteGroupRequest) (*charon.DeleteGroupResponse, error) {
	h := &deleteGroupHandler{
		handler: newHandler(rs, ctx, "delete_group"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "group has been deleted")

	return resp, err
}

// GetGroup implements charon.RPCServer interface.
func (rs *rpcServer) GetGroup(ctx context.Context, req *charon.GetGroupRequest) (*charon.GetGroupResponse, error) {
	h := &getGroupHandler{
		handler: newHandler(rs, ctx, "get_group"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "group has been retrieved")

	return resp, err
}

// ListGroups implements charon.RPCServer interface.
func (rs *rpcServer) ListGroups(ctx context.Context, req *charon.ListGroupsRequest) (*charon.ListGroupsResponse, error) {
	h := &listGroupsHandler{
		handler: newHandler(rs, ctx, "list_groups"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "list of groups has been retrieved")

	return resp, err
}

// ListGroupPermissions implements charon.RPCServer interface.
func (rs *rpcServer) ListGroupPermissions(ctx context.Context, req *charon.ListGroupPermissionsRequest) (*charon.ListGroupPermissionsResponse, error) {
	h := &listGroupPermissionsHandler{
		handler: newHandler(rs, ctx, "list_group_permissions"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "list of group permissions has been retrieved")

	return resp, err
}

// SetGroupPermissions implements charon.RPCServer interface.
func (rs *rpcServer) SetGroupPermissions(ctx context.Context, req *charon.SetGroupPermissionsRequest) (*charon.SetGroupPermissionsResponse, error) {
	h := &setGroupPermissionsHandler{
		handler: newHandler(rs, ctx, "set_group_permissions"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "group permissions has been set")

	return resp, err
}

// GetPermission implements charon.RPCServer interface.
func (rs *rpcServer) GetPermission(ctx context.Context, req *charon.GetPermissionRequest) (*charon.GetPermissionResponse, error) {
	h := &getPermissionHandler{
		handler: newHandler(rs, ctx, "get_permission"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "permission has been retrieved")

	return resp, err
}

// RegisterPermissions implements charon.RPCServer interface.
func (rs *rpcServer) RegisterPermissions(ctx context.Context, req *charon.RegisterPermissionsRequest) (*charon.RegisterPermissionsResponse, error) {
	h := &registerPermissionsHandler{
		handler:  newHandler(rs, ctx, "register_permissions"),
		registry: rs.permissionRegistry,
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "permissions has been registered")

	return resp, err
}

// ListPermissions implements charon.RPCServer interface.
func (rs *rpcServer) ListPermissions(ctx context.Context, req *charon.ListPermissionsRequest) (*charon.ListPermissionsResponse, error) {
	h := &listPermissionsHandler{
		handler: newHandler(rs, ctx, "list_permissions"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "list of permissions has been retrieved")

	return resp, err
}

// CreateUser implements charon.RPCServer interface.
func (rs *rpcServer) CreateUser(ctx context.Context, req *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	h := &createUserHandler{
		handler: newHandler(rs, ctx, "create_user"),
		hasher:  rs.passwordHasher,
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "user has been created")

	return resp, err
}

// ModifyUser implements charon.RPCServer interface.
func (rs *rpcServer) ModifyUser(ctx context.Context, req *charon.ModifyUserRequest) (*charon.ModifyUserResponse, error) {
	h := &modifyUserHandler{
		handler: newHandler(rs, ctx, "modify_user"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "user has been modified")

	return resp, err
}

// GetUser implements charon.RPCServer interface.
func (rs *rpcServer) GetUser(ctx context.Context, req *charon.GetUserRequest) (*charon.GetUserResponse, error) {
	h := &getUserHandler{
		handler: newHandler(rs, ctx, "get_user"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "user has been retrieved")

	return resp, err
}

// ListUsers implements charon.RPCServer interface.
func (rs *rpcServer) ListUsers(ctx context.Context, req *charon.ListUsersRequest) (*charon.ListUsersResponse, error) {
	h := &listUsersHandler{
		handler: newHandler(rs, ctx, "list_users"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "list of users has been retrieved")

	return resp, err
}

// DeleteUser implements charon.RPCServer interface.
func (rs *rpcServer) DeleteUser(ctx context.Context, req *charon.DeleteUserRequest) (*charon.DeleteUserResponse, error) {
	h := &deleteUserHandler{
		handler: newHandler(rs, ctx, "delete_user"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "user has been deleted")

	return resp, err
}

// SetUserGroups implements charon.RPCServer interface.
func (rs *rpcServer) SetUserGroups(ctx context.Context, req *charon.SetUserGroupsRequest) (*charon.SetUserGroupsResponse, error) {
	h := &setUserGroupsHandler{
		handler: newHandler(rs, ctx, "set_user_groups"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "user groups has been set")

	return resp, err
}

// ListUserGroups implements charon.RPCServer interface.
func (rs *rpcServer) ListUserGroups(ctx context.Context, req *charon.ListUserGroupsRequest) (*charon.ListUserGroupsResponse, error) {
	h := &listUserGroupsHandler{
		handler: newHandler(rs, ctx, "list_user_groups"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "list of user groups has been retrieved")

	return resp, err
}

// SetUserPermissions implements charon.RPCServer interface.
func (rs *rpcServer) SetUserPermissions(ctx context.Context, req *charon.SetUserPermissionsRequest) (*charon.SetUserPermissionsResponse, error) {
	h := &setUserPermissionsHandler{
		handler: newHandler(rs, ctx, "set_user_permissions"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "user permissions has been set")

	return resp, err
}

// ListUserPermissions implements charon.RPCServer interface.
func (rs *rpcServer) ListUserPermissions(ctx context.Context, req *charon.ListUserPermissionsRequest) (*charon.ListUserPermissionsResponse, error) {
	h := &listUserPermissionsHandler{
		handler: newHandler(rs, ctx, "list_user_permissions"),
	}
	h.monitor.requests.Add(1)

	resp, err := h.handle(ctx, req)
	h.handler.handle(err, "list of user permissions has been retrieved")

	return resp, err
}
