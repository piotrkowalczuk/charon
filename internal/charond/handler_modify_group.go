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
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc/codes"
)

type modifyGroupHandler struct {
	*handler
}

func (mgh *modifyGroupHandler) Modify(ctx context.Context, req *charonrpc.ModifyGroupRequest) (*charonrpc.ModifyGroupResponse, error) {
	if req.Id <= 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "group id is missing")
	}
	if !req.GetName().GetValid() && !req.GetDescription().GetValid() {
		return nil, grpcerr.E(codes.InvalidArgument, "nothing to be modified")
	}
	act, err := mgh.Actor(ctx)
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
			return nil, grpcerr.E(codes.NotFound, "group does not exists")
		}
		return nil, grpcerr.E(codes.Internal, "update group by id query failed", err)
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

	return grpcerr.E(codes.PermissionDenied, "group cannot be modified, missing permission")
}

func (mgh *modifyGroupHandler) response(ent *model.GroupEntity) (*charonrpc.ModifyGroupResponse, error) {
	msg, err := mapping.ReverseGroup(ent)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "group reverse mapping failure")
	}
	return &charonrpc.ModifyGroupResponse{Group: msg}, nil
}
