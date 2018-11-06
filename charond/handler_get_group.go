package charond

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getGroupHandler struct {
	*handler
}

func (ggh *getGroupHandler) Get(ctx context.Context, req *charonrpc.GetGroupRequest) (*charonrpc.GetGroupResponse, error) {
	act, err := ggh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = ggh.firewall(req, act); err != nil {
		return nil, err
	}

	ent, err := ggh.repository.group.FindOneByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "group does not exists")
		}
		return nil, grpc.Errorf(codes.Internal, err.Error())
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

	return grpc.Errorf(codes.PermissionDenied, "group cannot be retrieved, missing permission")
}

func (ggh *getGroupHandler) response(ent *model.GroupEntity) (*charonrpc.GetGroupResponse, error) {
	msg, err := ent.Message()
	if err != nil {
		return nil, err
	}
	return &charonrpc.GetGroupResponse{
		Group: msg,
	}, nil
}
