package main

import (
	"fmt"

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

	user, _, permissions, err := rs.retrieveUserData(ctx, token)
	if err != nil {
		return nil, err
	}

	fmt.Println(user, permissions)
	if r.IsSuperuser != nil && r.IsSuperuser.Valid && !user.IsSuperuser && !permissions.Contains(charon.UserCanCreateSuper) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create user with is_superuser property that has custom value")
	}

	if r.IsStaff != nil && r.IsStaff.Valid && !user.IsSuperuser && !permissions.Contains(charon.UserCanCreateStaff) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create user with is_staff property that has custom value")
	}

	if r.IsActive != nil && r.IsActive.Valid && !user.IsSuperuser && !permissions.Contains(charon.UserCanCreateActive) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create user with is_active property that has custom value")
	}

	if r.IsConfirmed != nil && r.IsConfirmed.Valid && !user.IsSuperuser && !permissions.Contains(charon.UserCanCreateConfirmed) {
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
		return nil, mapUserError(err)
	}

	return &charon.CreateUserResponse{
		User: entity.Message(),
	}, nil
}

// ModifyUser ...
func (rs *rpcServer) ModifyUser(ctx context.Context, r *charon.ModifyUserRequest) (*charon.ModifyUserResponse, error) {
	user, err := rs.userRepository.UpdateOneByID(
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
		User: user.Message(),
	}, nil
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
