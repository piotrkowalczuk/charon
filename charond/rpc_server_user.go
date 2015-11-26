package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/protot"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateUser ...
func (rs *rpcServer) CreateUser(ctx context.Context, r *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	token, err := rs.token(ctx)
	if err != nil {
		return nil, err
	}
	user, _, permissions, err := rs.retrieveUserData(ctx, token)
	if err != nil {
		return nil, err
	}

	if r.IsSuperuser.Valid && !user.IsSuperuser && !permissions.Contains(charon.UserCanCreateSuper) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create user with is_superuser property that has custom value")
	}

	if r.IsStaff.Valid && !user.IsSuperuser && !permissions.Contains(charon.UserCanCreateStaff) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create user with is_staff property that has custom value")
	}

	if r.IsActive.Valid && !user.IsSuperuser && !permissions.Contains(charon.UserCanCreateActive) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create user with is_active property that has custom value")
	}

	if r.IsConfirmed.Valid && !user.IsSuperuser && !permissions.Contains(charon.UserCanCreateConfirmed) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create user with is_confirmed property that has custom value")
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
		return nil, err
	}

	return &charon.CreateUserResponse{
		Id:        entity.ID,
		CreatedAt: protot.TimeToTimestamp(*entity.CreatedAt),
	}, nil
}

// ModifyUser ...
func (rs *rpcServer) ModifyUser(ctx context.Context, r *charon.ModifyUserRequest) (*charon.ModifyUserResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "modify user is not implemented yet")
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

// ModifyUserPassword ...
func (rs *rpcServer) ModifyUserPassword(ctx context.Context, r *charon.ModifyUserPasswordRequest) (*charon.ModifyUserPasswordResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "modify user password is not implemented yet")
}
