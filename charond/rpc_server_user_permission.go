package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
)

// ListUserGroups implements charon.RPCServer interface.
func (rs *rpcServer) ListUserGroups(ctx context.Context, req *charon.ListUserGroupsRequest) (*charon.ListUserGroupsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: list user groups endpoint is not implemented yet")
}

// SetUserGroups implements charon.RPCServer interface.
func (rs *rpcServer) SetUserGroups(ctx context.Context, req *charon.SetUserGroupsRequest) (*charon.SetUserGroupsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: set user groups endpoint is not implemented yet")
}
