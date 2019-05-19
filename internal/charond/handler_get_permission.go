package charond

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/session"

	"google.golang.org/grpc/codes"
)

type getPermissionHandler struct {
	*handler
}

func (gph *getPermissionHandler) Get(ctx context.Context, req *charonrpc.GetPermissionRequest) (*charonrpc.GetPermissionResponse, error) {
	if req.Id < 1 {
		return nil, grpcerr.E(codes.InvalidArgument, "permission id needs to be greater than zero")
	}

	act, err := gph.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = gph.firewall(req, act); err != nil {
		return nil, err
	}

	permission, err := gph.repository.permission.FindOneByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpcerr.E(codes.NotFound, "permission does not exists")
		}
		return nil, grpcerr.E(codes.Internal, "permission cannot be fetched", err)
	}

	return &charonrpc.GetPermissionResponse{
		Permission: permission.Permission().String(),
	}, nil
}

func (gph *getPermissionHandler) firewall(req *charonrpc.GetPermissionRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.PermissionCanRetrieve) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "permission cannot be retrieved, missing permission")
}
