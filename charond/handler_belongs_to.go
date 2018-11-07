package charond

import (
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/session"
	"google.golang.org/grpc/codes"
)

type belongsToHandler struct {
	*handler
}

func (bth *belongsToHandler) BelongsTo(ctx context.Context, req *charonrpc.BelongsToRequest) (*wrappers.BoolValue, error) {
	if req.GroupId < 1 {
		return nil, grpcerr.E(codes.InvalidArgument, "group id needs to be greater than zero")
	}
	if req.UserId < 1 {
		return nil, grpcerr.E(codes.InvalidArgument, "user id needs to be greater than zero")
	}

	act, err := bth.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = bth.firewall(req, act); err != nil {
		return nil, err
	}

	belongs, err := bth.repository.userGroups.Exists(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "user group fetch failure", err)
	}

	return &wrappers.BoolValue{Value: belongs}, nil
}

func (bth *belongsToHandler) firewall(req *charonrpc.BelongsToRequest, act *session.Actor) error {
	if act.User.ID == req.UserId {
		return nil
	}
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.UserGroupCanCheckBelongingAsStranger) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "group belonging cannot be checked, missing permission")
}
