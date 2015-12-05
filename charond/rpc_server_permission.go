package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// RegisterPermissions implements charon.RPCServer interface.
func (rs *rpcServer) RegisterPermissions(ctx context.Context, req *charon.RegisterPermissionsRequest) (*charon.RegisterPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: registrater permissions endpoint is not implemented yet")
}

// ListPermissions implements charon.RPCServer interface.
func (rs *rpcServer) ListPermissions(ctx context.Context, req *charon.ListPermissionsRequest) (*charon.ListPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: list permissions endpoint is not implemented yet")
}
