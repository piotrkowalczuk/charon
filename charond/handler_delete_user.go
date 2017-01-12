package charond

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type deleteUserHandler struct {
	*handler
}

func (duh *deleteUserHandler) Delete(ctx context.Context, req *charonrpc.DeleteUserRequest) (*wrappers.BoolValue, error) {
	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "User cannot be deleted, invalid ID: %d", req.Id)
	}

	act, err := duh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	ent, err := duh.repository.user.FindOneByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if err = duh.firewall(req, act, ent); err != nil {
		return nil, err
	}

	affected, err := duh.repository.user.DeleteOneByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "User does not exists")
		}
		return nil, err
	}

	return &wrappers.BoolValue{
		Value: affected > 0,
	}, nil
}

func (duh *deleteUserHandler) firewall(req *charonrpc.DeleteUserRequest, act *session.Actor, ent *model.UserEntity) error {
	if act.User.ID == ent.ID {
		return grpc.Errorf(codes.PermissionDenied, "User is not permited to remove himself")
	}
	if act.User.IsSuperuser {
		return nil
	}
	if ent.IsSuperuser {
		return grpc.Errorf(codes.PermissionDenied, "only superuser can remove other superuser")
	}
	if ent.IsStaff {
		switch {
		case act.User.ID == ent.CreatedBy.Int64Or(0):
			if !act.Permissions.Contains(charon.UserCanDeleteStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "staff User cannot be removed by owner, missing permission")
			}
			return nil
		case !act.Permissions.Contains(charon.UserCanDeleteStaffAsStranger):
			return grpc.Errorf(codes.PermissionDenied, "staff User cannot be removed by stranger, missing permission")
		}
		return nil
	}

	if act.User.ID == ent.CreatedBy.Int64Or(0) {
		if !act.Permissions.Contains(charon.UserCanDeleteAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "User cannot be removed by owner, missing permission")
		}
		return nil
	}
	if !act.Permissions.Contains(charon.UserCanDeleteAsStranger) {
		return grpc.Errorf(codes.PermissionDenied, "User cannot be removed by stranger, missing permission")
	}
	return nil
}
