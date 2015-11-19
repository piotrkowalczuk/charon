package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// GetUserPermissions implements charon.RPCServer interface.
func (rs *rpcServer) GetUserPermissions(ctx context.Context, req *charon.GetUserPermissionsRequest) (*charon.GetUserPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "get user permissions is not implemented yet")
}
