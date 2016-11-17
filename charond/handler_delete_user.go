package charond

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type deleteUserHandler struct {
	*handler
}

func (duh *deleteUserHandler) Delete(ctx context.Context, req *charonrpc.DeleteUserRequest) (*wrappers.BoolValue, error) {
	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "user cannot be deleted, invalid ID: %d", req.Id)
	}

	act, err := duh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	ent, err := duh.repository.user.FindOneByID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = duh.firewall(req, act, ent); err != nil {
		return nil, err
	}

	affected, err := duh.repository.user.DeleteOneByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "user does not exists")
		}
		return nil, err
	}

	return &wrappers.BoolValue{
		Value: affected > 0,
	}, nil
}

func (duh *deleteUserHandler) firewall(req *charonrpc.DeleteUserRequest, act *actor, ent *model.UserEntity) error {
	if act.user.ID == ent.ID {
		return grpc.Errorf(codes.PermissionDenied, "user is not permited to remove himself")
	}
	if act.user.IsSuperuser {
		return nil
	}
	if ent.IsSuperuser {
		return grpc.Errorf(codes.PermissionDenied, "only superuser can remove other superuser")
	}
	if ent.IsStaff {
		switch {
		case act.user.ID == ent.CreatedBy.Int64Or(0):
			if !act.permissions.Contains(charon.UserCanDeleteStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "staff user cannot be removed by owner, missing permission")
			}
			return nil
		case !act.permissions.Contains(charon.UserCanDeleteStaffAsStranger):
			return grpc.Errorf(codes.PermissionDenied, "staff user cannot be removed by stranger, missing permission")
		}
		return nil
	}

	if act.user.ID == ent.CreatedBy.Int64Or(0) {
		if !act.permissions.Contains(charon.UserCanDeleteAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "user cannot be removed by owner, missing permission")
		}
		return nil
	}
	if !act.permissions.Contains(charon.UserCanDeleteAsStranger) {
		return grpc.Errorf(codes.PermissionDenied, "user cannot be removed by stranger, missing permission")
	}
	return nil
}
