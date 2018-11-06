package charond

import (
	"context"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"

	"google.golang.org/grpc/codes"
)

type listGroupsHandler struct {
	*handler
}

func (lgh *listGroupsHandler) List(ctx context.Context, req *charonrpc.ListGroupsRequest) (*charonrpc.ListGroupsResponse, error) {
	act, err := lgh.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = lgh.firewall(req, act); err != nil {
		return nil, err
	}

	ents, err := lgh.repository.group.Find(ctx, &model.GroupFindExpr{
		Limit:   req.GetLimit().Int64Or(10),
		Offset:  req.GetOffset().Int64Or(0),
		OrderBy: mapping.OrderBy(req.GetOrderBy()),
	})
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "find group query failed", err)
	}

	return lgh.response(ents)
}

func (lgh *listGroupsHandler) firewall(req *charonrpc.ListGroupsRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupCanRetrieve) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "list of groups cannot be retrieved, missing permission")
}

func (lgh *listGroupsHandler) response(ents []*model.GroupEntity) (*charonrpc.ListGroupsResponse, error) {
	msg, err := mapping.ReverseGroups(ents)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "group entities mapping failure", err)

	}

	return &charonrpc.ListGroupsResponse{Groups: msg}, nil
}
