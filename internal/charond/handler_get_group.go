package charond

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"google.golang.org/grpc/codes"
)

type getGroupHandler struct {
	*handler
}

func (ggh *getGroupHandler) Get(ctx context.Context, req *charonrpc.GetGroupRequest) (*charonrpc.GetGroupResponse, error) {
	if req.Id <= 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "missing group id")
	}
	act, err := ggh.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = ggh.firewall(req, act); err != nil {
		return nil, err
	}

	ent, err := ggh.repository.group.FindOneByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpcerr.E(codes.NotFound, "group does not exists")
		}
		return nil, grpcerr.E(codes.Internal, "group cannot be fetched", err)
	}

	return ggh.response(ent)
}

func (ggh *getGroupHandler) firewall(req *charonrpc.GetGroupRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupCanRetrieve) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "group cannot be retrieved, missing permission")
}

func (ggh *getGroupHandler) response(ent *model.GroupEntity) (*charonrpc.GetGroupResponse, error) {
	msg, err := mapping.ReverseGroup(ent)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "group entity mapping failure", err)
	}
	return &charonrpc.GetGroupResponse{
		Group: msg,
	}, nil
}
