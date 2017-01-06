package charond

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/session"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type deleteGroupHandler struct {
	*handler
}

func (dgh *deleteGroupHandler) Delete(ctx context.Context, req *charonrpc.DeleteGroupRequest) (*wrappers.BoolValue, error) {
	act, err := dgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = dgh.firewall(req, act); err != nil {
		return nil, err
	}

	affected, err := dgh.repository.group.DeleteOneByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "group does not exists")
		}
		return nil, err
	}

	return &wrappers.BoolValue{
		Value: affected > 0,
	}, nil
}

func (dgh *deleteGroupHandler) firewall(req *charonrpc.DeleteGroupRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupCanDelete) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "group cannot be removed, missing permission")
}
