package charond

import (
	"context"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/session"

	"google.golang.org/grpc/codes"
)

type listUserGroupsHandler struct {
	*handler
}

func (lugh *listUserGroupsHandler) ListGroups(ctx context.Context, req *charonrpc.ListUserGroupsRequest) (*charonrpc.ListUserGroupsResponse, error) {
	if req.Id <= 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "missing user id")
	}
	act, err := lugh.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = lugh.firewall(req, act); err != nil {
		return nil, err
	}

	ents, err := lugh.repository.group.FindByUserID(ctx, req.Id)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "find groups by user id query failed", err)
	}

	msg, err := mapping.ReverseGroups(ents)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "user group entities mapping failure", err)
	}

	return &charonrpc.ListUserGroupsResponse{Groups: msg}, nil
}

func (lugh *listUserGroupsHandler) firewall(req *charonrpc.ListUserGroupsRequest, act *session.Actor) error {
	switch {
	case act.User.IsSuperuser:
		return nil
	case act.User.ID == req.Id:
		return nil
	case act.Permissions.Contains(charon.UserGroupCanRetrieve):
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "list of user groups cannot be retrieved, missing permission")
}
