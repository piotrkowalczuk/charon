package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ListGroupPermissions implements charon.RPCServer interface.
func (rs *rpcServer) ListGroupPermissions(ctx context.Context, req *charon.ListGroupPermissionsRequest) (*charon.ListGroupPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: list group permissions endpoint is not implemented yet")
}

// SetGroupPermissions implements charon.RPCServer interface.
func (rs *rpcServer) SetGroupPermissions(ctx context.Context, req *charon.SetGroupPermissionsRequest) (*charon.SetGroupPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: set group permissions endpoint is not implemented yet")
}
