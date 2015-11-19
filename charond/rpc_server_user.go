package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateUser ...
func (rs *rpcServer) CreateUser(ctx context.Context, r *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	session, err := rs.mnemosyne.GetArbitrarily(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := charon.UserIDFromSession(session)
	if err != nil {
		return nil, err
	}

	// TODO: take into account permissions (do not allow to create superuser/staff user without necessary permissions)
	_, err = rs.permissionRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if r.SecurePassword == "" {
		r.SecurePassword, err = rs.passwordHasher.Hash(r.PlainPassword)
		if err != nil {
			return nil, err
		}
	}

	entity, err := rs.userRepository.Create(r.Username, r.SecurePassword, r.FirstName, r.LastName, uuid.New())
	if err != nil {
		return nil, err
	}

	return &charon.CreateUserResponse{
		Id:        entity.ID,
		CreatedAt: mnemosyne.TimeToTimestamp(*entity.CreatedAt),
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
