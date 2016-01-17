package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
)

// SetUserGroups implements charon.RPCServer interface.
func (rs *rpcServer) SetUserGroups(ctx context.Context, req *charon.SetUserGroupsRequest) (*charon.SetUserGroupsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: set user groups endpoint is not implemented yet")
}

// ListUserGroups implements charon.RPCServer interface.
func (rs *rpcServer) ListUserGroups(ctx context.Context, req *charon.ListUserGroupsRequest) (*charon.ListUserGroupsResponse, error) {
	entities, err := rs.repository.group.FindByUserID(req.Id)
	if err != nil {
		return nil, err
	}

	groups := make([]*charon.Group, 0, len(entities))
	for _, e := range entities {
		groups = append(groups, e.Message())
	}

	return &charon.ListUserGroupsResponse{
		Groups: groups,
	}, nil
}
