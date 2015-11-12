package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateUser ...
func (rs *rpcServer) CreateUser(ctx context.Context, r *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "create user is not implemented yet")
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
