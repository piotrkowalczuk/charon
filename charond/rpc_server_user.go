package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/protot"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateUser ...
func (rs *rpcServer) CreateUser(ctx context.Context, r *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	user, _, permissions, err := rs.retrieveUserData(ctx)
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
	return nil, grpc.Errorf(codes.Unimplemented, "get user is not implemented yet")
}

// GetUsers ...
func (rs *rpcServer) GetUsers(ctx context.Context, r *charon.GetUsersRequest) (*charon.GetUsersResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "get users is not implemented yet")
}

// DeleteUser ...
func (rs *rpcServer) DeleteUser(ctx context.Context, r *charon.DeleteUserRequest) (*charon.DeleteUserResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "delete user is not implemented yet")
}

// ModifyUserPassword ...
func (rs *rpcServer) ModifyUserPassword(ctx context.Context, r *charon.ModifyUserPasswordRequest) (*charon.ModifyUserPasswordResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "modify user password is not implemented yet")
}
