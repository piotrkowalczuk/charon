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

// CreateUser ...
func (rs *rpcServer) CreateUser(ctx context.Context, r *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	var err error
	defer func() {
		if err != nil {
			sklog.Error(rs.logger, err)
		} else {
			sklog.Debug(rs.logger, "user created")
		}
	}()

	token, err := rs.token(ctx)
	if err != nil {
		return nil, err
	}

	user, _, permissions, err := rs.retrieveActor(ctx, token)
	if err != nil {
		return nil, err
	}

	if !user.IsSuperuser {
		if r.IsSuperuser != nil && r.IsSuperuser.Valid {
			return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create superuser")
		}

		if r.IsStaff != nil && r.IsStaff.Valid && !permissions.Contains(charon.UserCanCreateStaff) {
			return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create staff user")
		}
	}

	if r.SecurePassword == "" {
		r.SecurePassword, err = rs.passwordHasher.Hash(r.PlainPassword)
		if err != nil {
			return nil, err
		}
	} else {
		if !user.IsSuperuser {
			return nil, grpc.Errorf(codes.PermissionDenied, "charond: only superuser can create an user with manualy defined secure password")
		}
	}

	entity, err := rs.userRepository.Create(
		r.Username,
		r.SecurePassword,
		r.FirstName,
		r.LastName,
		uuid.New(),
		r.IsSuperuser.BoolOr(false),
		r.IsStaff.BoolOr(false),
		r.IsActive.BoolOr(false),
		r.IsConfirmed.BoolOr(false),
	)
	if err != nil {
		return nil, mapUserError(err)
	}

	return &charon.CreateUserResponse{
		User: entity.Message(),
	}, nil
}

// ModifyUser ...
func (rs *rpcServer) ModifyUser(ctx context.Context, r *charon.ModifyUserRequest) (*charon.ModifyUserResponse, error) {
	token, err := rs.token(ctx)
	if err != nil {
		return nil, err
	}

	actor, _, permissions, err := rs.retrieveActor(ctx, token)
	if err != nil {
		return nil, err
	}

	entity, err := rs.userRepository.FindOneByID(r.Id)
	if err != nil {
		return nil, err
	}

	if hint, ok := modifyUserFirewall(r, entity, actor, permissions); !ok {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: "+hint)
	}

	entity, err = rs.userRepository.UpdateOneByID(
		r.Id,
		r.Username,
		r.SecurePassword,
		r.FirstName,
		r.LastName,
		r.IsSuperuser,
		r.IsActive,
		r.IsStaff,
		r.IsConfirmed,
	)
	if err != nil {
		return nil, mapUserError(err)
	}

	sklog.Debug(rs.logger, "user modified", "id", r.Id)

	return &charon.ModifyUserResponse{
		User: entity.Message(),
	}, nil
}

func modifyUserFirewall(r *charon.ModifyUserRequest, entity *userEntity, actor *userEntity, perms charon.Permissions) (string, bool) {
	isOwner := actor.ID == entity.ID

	if !actor.IsSuperuser {
		switch {
		case entity.IsSuperuser:
			return "only superuser can modify a superuser account", false
		case entity.IsStaff && !isOwner && perms.Contains(charon.UserCanModifyStaffAsStranger):
			return "missing permission to modify an account as a stranger", false
		case entity.IsStaff && isOwner && perms.Contains(charon.UserCanModifyStaffAsOwner):
			return "missing permission to modify an account as an owner", false
		case r.IsSuperuser != nil && r.IsSuperuser.Valid:
			return "only superuser can change existing account to superuser", false
		case r.IsStaff != nil && r.IsStaff.Valid && !perms.Contains(charon.UserCanCreateStaff):
			return "user is not allowed to create user with is_staff property that has custom value", false
		}
	}

	return "", true
}

// GetUser ...
func (rs *rpcServer) GetUser(ctx context.Context, r *charon.GetUserRequest) (*charon.GetUserResponse, error) {
	user, err := rs.userRepository.FindOneByID(r.Id)
	if err != nil {
		return nil, err
	}

	sklog.Debug(rs.logger, "user retrieved", "id", r.Id)

	return &charon.GetUserResponse{
		User: user.Message(),
	}, nil
}

// GetUsers ...
func (rs *rpcServer) GetUsers(ctx context.Context, r *charon.GetUsersRequest) (*charon.GetUsersResponse, error) {
	users, err := rs.userRepository.Find(r.Offset, r.Limit)
	if err != nil {
		return nil, err
	}

	resp := &charon.GetUsersResponse{
		Users: make([]*charon.User, 0, len(users)),
	}

	for _, u := range users {
		resp.Users = append(resp.Users, u.Message())
	}

	sklog.Debug(rs.logger, "users list retrieved", "count", len(users))

	return resp, nil
}

// DeleteUser ...
func (rs *rpcServer) DeleteUser(ctx context.Context, r *charon.DeleteUserRequest) (*charon.DeleteUserResponse, error) {
	if r.Id == 0 {
		return nil, grpc.Errorf(codes.FailedPrecondition, "charond: user id needs to be greater than zero")
	}
	affected, err := rs.userRepository.DeleteOneByID(r.Id)
	if err != nil {
		return nil, err
	}

	sklog.Debug(rs.logger, "users deleted", "id", r.Id)

	return &charon.DeleteUserResponse{
		Affected: affected,
	}, nil
}

func mapUserError(err error) error {
	switch pqcnstr.FromError(err) {
	case sqlCnstrPrimaryKeyUser:
		return grpc.Errorf(codes.AlreadyExists, charon.ErrDescUserWithIDExists)
	case sqlCnstrUniqueUserUsername:
		return grpc.Errorf(codes.AlreadyExists, charon.ErrDescUserWithUsernameExists)
	default:
		return err
	}
}
