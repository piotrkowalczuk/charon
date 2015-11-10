package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
)

// Login ...
func (rs *rpcServer) Login(ctx context.Context, r *charon.LoginRequest) (*charon.LoginResponse, error) {
	return nil, nil
}

// Logout ...
func (rs *rpcServer) Logout(ctx context.Context, r *charon.LogoutRequest) (*charon.LogoutResponse, error) {
	return nil, nil
}

// IsGranted ...
func (rs *rpcServer) IsGranted(ctx context.Context, r *charon.IsGrantedRequest) (*charon.IsGrantedResponse, error) {
	return nil, nil
}

// HasPrivileges ...
func (rs *rpcServer) HasPrivileges(ctx context.Context, r *charon.HasPrivilegesRequest) (*charon.HasPrivilegesResponse, error) {
	return nil, nil
}
