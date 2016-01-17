package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/pqcnstr"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateUser implements charon.RPCServer interface.
func (rs *rpcServer) CreateUser(ctx context.Context, req *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	var err error
	defer func() {
		if err != nil {
			sklog.Error(rs.logger, err)
		} else {
			sklog.Debug(rs.logger, "user created", "username", req.Username, "superuser", req.IsSuperuser.BoolOr(false))
		}
	}()

	actor, err := rs.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	if !actor.isLocalhost() && !actor.user.IsSuperuser {
		if req.IsSuperuser.BoolOr(false) {
			return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create superuser")
		}

		if req.IsStaff.BoolOr(false) && !actor.permissions.Contains(charon.UserCanCreateStaff) {
			return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create staff user")
		}
	}

	if len(req.SecurePassword) == 0 {
		req.SecurePassword, err = rs.passwordHasher.Hash([]byte(req.PlainPassword))
		if err != nil {
			return nil, err
		}
	} else {
		if !actor.user.IsSuperuser {
			return nil, grpc.Errorf(codes.PermissionDenied, "charond: only superuser can create an user with manualy defined secure password")
		}
	}

	entity, err := rs.repository.user.Create(
		req.Username,
		req.SecurePassword,
		req.FirstName,
		req.LastName,
		uuid.NewRandom(),
		req.IsSuperuser.BoolOr(false),
		req.IsStaff.BoolOr(false),
		req.IsActive.BoolOr(false),
		req.IsConfirmed.BoolOr(false),
	)
	if err != nil {
		return nil, mapUserError(err)
	}

	return &charon.CreateUserResponse{
		User: entity.Message(),
	}, nil
}

// ModifyUser implements charon.RPCServer interface.
func (rs *rpcServer) ModifyUser(ctx context.Context, req *charon.ModifyUserRequest) (*charon.ModifyUserResponse, error) {
	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: user cannot be modified, invalid id: %d", req.Id)
	}

	actor, err := rs.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := rs.repository.user.FindOneByID(req.Id)
	if err != nil {
		return nil, err
	}

	if hint, ok := modifyUserFirewall(req, entity, actor); !ok {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: "+hint)
	}

	entity, err = rs.repository.user.UpdateOneByID(
		req.Id,
		req.Username,
		req.SecurePassword,
		req.FirstName,
		req.LastName,
		req.IsSuperuser,
		req.IsActive,
		req.IsStaff,
		req.IsConfirmed,
	)
	if err != nil {
		return nil, mapUserError(err)
	}

	sklog.Debug(rs.logger, "user modified", "id", req.Id)

	return &charon.ModifyUserResponse{
		User: entity.Message(),
	}, nil
}

func modifyUserFirewall(req *charon.ModifyUserRequest, entity *userEntity, actor *actor) (string, bool) {
	isOwner := actor.user.ID == entity.ID

	if !actor.user.IsSuperuser {
		switch {
		case entity.IsSuperuser:
			return "only superuser can modify a superuser account", false
		case entity.IsStaff && !isOwner && actor.permissions.Contains(charon.UserCanModifyStaffAsStranger):
			return "missing permission to modify an account as a stranger", false
		case entity.IsStaff && isOwner && actor.permissions.Contains(charon.UserCanModifyStaffAsOwner):
			return "missing permission to modify an account as an owner", false
		case req.IsSuperuser != nil && req.IsSuperuser.Valid:
			return "only superuser can change existing account to superuser", false
		case req.IsStaff != nil && req.IsStaff.Valid && !actor.permissions.Contains(charon.UserCanCreateStaff):
			return "user is not allowed to create user with is_staff property that has custom value", false
		}
	}

	return "", true
}

// GetUser implements charon.RPCServer interface.
func (rs *rpcServer) GetUser(ctx context.Context, req *charon.GetUserRequest) (*charon.GetUserResponse, error) {
	user, err := rs.repository.user.FindOneByID(req.Id)
	if err != nil {
		return nil, err
	}

	sklog.Debug(rs.logger, "user retrieved", "id", req.Id)

	return &charon.GetUserResponse{
		User: user.Message(),
	}, nil
}

// ListUsers implements charon.RPCServer interface.
func (rs *rpcServer) ListUsers(ctx context.Context, req *charon.ListUsersRequest) (*charon.ListUsersResponse, error) {
	users, err := rs.repository.user.Find(req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	resp := &charon.ListUsersResponse{
		Users: make([]*charon.User, 0, len(users)),
	}

	for _, u := range users {
		resp.Users = append(resp.Users, u.Message())
	}

	sklog.Debug(rs.logger, "users list retrieved", "count", len(users))

	return resp, nil
}

// DeleteUser implements charon.RPCServer interface.
func (rs *rpcServer) DeleteUser(ctx context.Context, req *charon.DeleteUserRequest) (*charon.DeleteUserResponse, error) {
	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: user cannot be deleted, invalid id: %d", req.Id)
	}
	affected, err := rs.repository.user.DeleteOneByID(req.Id)
	if err != nil {
		return nil, err
	}

	sklog.Debug(rs.logger, "users deleted", "id", req.Id)

	return &charon.DeleteUserResponse{
		Affected: affected,
	}, nil
}

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
