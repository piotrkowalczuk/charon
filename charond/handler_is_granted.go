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

type isGrantedHandler struct {
	*handler
}

func (ig *isGrantedHandler) IsGranted(ctx context.Context, req *charonrpc.IsGrantedRequest) (*wrappers.BoolValue, error) {
	if req.Permission == "" {
		return nil, grpcerr.E(codes.InvalidArgument, "permission cannot be empty")
	}
	if req.UserId < 1 {
		return nil, grpcerr.E(codes.InvalidArgument, "user id needs to be greater than zero")
	}

	act, err := ig.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = ig.firewall(req, act); err != nil {
		return nil, err
	}

	granted, err := ig.repository.user.IsGranted(ctx, req.UserId, charon.Permission(req.Permission))
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "is granted repository call failure", err)
	}

	return &wrappers.BoolValue{Value: granted}, nil
}

func (ig *isGrantedHandler) firewall(req *charonrpc.IsGrantedRequest, act *session.Actor) error {
	if act.User.ID == req.UserId {
		return nil
	}
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.UserPermissionCanCheckGrantingAsStranger) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "group granting cannot be checked, missing permission")
}
