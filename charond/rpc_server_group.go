package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateGroup implements charon.RPCServer interface.
func (rs *rpcServer) CreateGroup(ctx context.Context, req *charon.CreateGroupRequest) (*charon.CreateGroupResponse, error) {
	actor, err := rs.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	group, err := rs.repository.group.Create(req.Name, req.Description, actor.user.ID)
	if err != nil {
		return nil, err
	}

	return &charon.CreateGroupResponse{
		Group: group.Message(),
	}, nil
}

// ModifyGroup implements charon.RPCServer interface.
func (rs *rpcServer) ModifyGroup(ctx context.Context, req *charon.ModifyGroupRequest) (*charon.ModifyGroupResponse, error) {
	actor, err := rs.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	group, err := rs.repository.group.UpdateOneByID(req.Id, actor.user.ID, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &charon.ModifyGroupResponse{
		Group: group.Message(),
	}, nil
}

// DeleteGroup implements charon.RPCServer interface.
func (rs *rpcServer) DeleteGroup(ctx context.Context, req *charon.DeleteGroupRequest) (*charon.DeleteGroupResponse, error) {
	affected, err := rs.repository.group.DeleteOneByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &charon.DeleteGroupResponse{
		Affected: affected,
	}, nil
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
