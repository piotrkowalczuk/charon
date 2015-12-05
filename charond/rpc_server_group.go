package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateGroup implements charon.RPCServer interface.
func (rs *rpcServer) CreateGroup(ctx context.Context, req *charon.CreateGroupRequest) (*charon.CreateGroupResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: create group endpoint is not implemented yet")
}

// ModifyGroup implements charon.RPCServer interface.
func (rs *rpcServer) ModifyGroup(ctx context.Context, req *charon.ModifyGroupRequest) (*charon.ModifyGroupResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: modify group endpoint is not implemented yet")
}

// DeleteGroup implements charon.RPCServer interface.
func (rs *rpcServer) DeleteGroup(ctx context.Context, req *charon.DeleteGroupRequest) (*charon.DeleteGroupResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: delete group endpoint is not implemented yet")
}

// GetGroup implements charon.RPCServer interface.
func (rs *rpcServer) GetGroup(ctx context.Context, req *charon.GetGroupRequest) (*charon.GetGroupResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: get group endpoint is not implemented yet")
}

// ListGroups implements charon.RPCServer interface.
func (rs *rpcServer) ListGroups(ctx context.Context, req *charon.ListGroupsRequest) (*charon.ListGroupsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: list groups endpoint is not implemented yet")
}

// ListGroupPermissions implements charon.RPCServer interface.
func (rs *rpcServer) ListGroupPermissions(ctx context.Context, req *charon.ListGroupPermissionsRequest) (*charon.ListGroupPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: list group permissions endpoint is not implemented yet")
}

// SetGroupPermissions implements charon.RPCServer interface.
func (rs *rpcServer) SetGroupPermissions(ctx context.Context, req *charon.SetGroupPermissionsRequest) (*charon.SetGroupPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: set group permissions endpoint is not implemented yet")
}
