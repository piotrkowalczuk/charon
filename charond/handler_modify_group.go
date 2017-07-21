package charond

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/ntypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type modifyGroupHandler struct {
	*handler
}

func (mgh *modifyGroupHandler) Modify(ctx context.Context, req *charonrpc.ModifyGroupRequest) (*charonrpc.ModifyGroupResponse, error) {
	act, err := mgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	if err = mgh.firewall(req, act); err != nil {
		return nil, err
	}

	group, err := mgh.repository.group.UpdateOneByID(ctx, req.Id, &model.GroupPatch{
		UpdatedBy:   ntypes.Int64{Int64: act.User.ID, Valid: true},
		Name:        allocNilString(req.Name),
		Description: allocNilString(req.Description),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "group does not exists")
		}
		return nil, err
	}

	return mgh.response(group)
}

func (mgh *modifyGroupHandler) firewall(req *charonrpc.ModifyGroupRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupCanModify) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "group cannot be modified, missing permission")
}

func (mgh *modifyGroupHandler) response(g *model.GroupEntity) (*charonrpc.ModifyGroupResponse, error) {
	msg, err := g.Message()
	if err != nil {
		return nil, err
	}
	return &charonrpc.ModifyGroupResponse{
		Group: msg,
	}, nil
}
