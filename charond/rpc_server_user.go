package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
)

// CreateUser ...
func (rs *rpcServer) CreateUser(ctx context.Context, r *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	return nil, nil
}

// ModifyUser ...
func (rs *rpcServer) ModifyUser(ctx context.Context, r *charon.ModifyUserRequest) (*charon.ModifyUserResponse, error) {
	return nil, nil
}

// GetUser ...
func (rs *rpcServer) GetUser(ctx context.Context, r *charon.GetUserRequest) (*charon.GetUserResponse, error) {
	return nil, nil
}

// GetUsers ...
func (rs *rpcServer) GetUsers(ctx context.Context, r *charon.GetUsersRequest) (*charon.GetUsersResponse, error) {
	return nil, nil
}

// DeleteUser ...
func (rs *rpcServer) DeleteUser(ctx context.Context, r *charon.DeleteUserRequest) (*charon.DeleteUserResponse, error) {
	return nil, nil
}
