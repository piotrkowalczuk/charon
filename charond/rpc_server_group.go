package main

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateGroup implements charon.RPCServer interface.
func (rs *rpcServer) CreateGroup(ctx context.Context, req *charon.CreateGroupRequest) (*charon.CreateGroupResponse, error) {
	token, err := rs.token(ctx)
	if err != nil {
		return nil, err
	}

	actor, _, permissions, err := rs.retrieveActor(ctx, token)
	if err != nil {
		return nil, err
	}

	if !permissions.Contains(charon.GroupCanCreate) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: actor do not have permission: %s", charon.GroupCanCreate.String())
	}

	entity, err := rs.repository.group.Create(actor.ID, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &charon.CreateGroupResponse{
		Group: entity.Message(),
	}, nil
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
	entity, err := rs.repository.group.FindOneByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "charond: group with id %d does not exists", req.Id)
		}
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &charon.GetGroupResponse{
		Group: entity.Message(),
	}, nil
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
