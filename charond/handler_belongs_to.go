package charond

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type belongsToHandler struct {
	*handler
}

func (bth *belongsToHandler) BelongsTo(ctx context.Context, req *charonrpc.BelongsToRequest) (*wrappers.BoolValue, error) {
	if req.GroupId < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "group id needs to be greater than zero")
	}
	if req.UserId < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "user id needs to be greater than zero")
	}

	act, err := bth.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = bth.firewall(req, act); err != nil {
		return nil, err
	}

	belongs, err := bth.repository.userGroups.Exists(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}

	return &wrappers.BoolValue{Value: belongs}, nil
}

func (bth *belongsToHandler) firewall(req *charonrpc.BelongsToRequest, act *actor) error {
	if act.user.ID == req.UserId {
		return nil
	}
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.UserGroupCanCheckBelongingAsStranger) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "group belonging cannot be checked, missing permission")
}
