package charond

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type deleteUserHandler struct {
	*handler
}

func (duh *deleteUserHandler) handle(ctx context.Context, req *charon.DeleteUserRequest) (*charon.DeleteUserResponse, error) {
	duh.loggerWith("user_id", req.Id)

	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "user cannot be deleted, invalid id: %d", req.Id)
	}

	act, err := duh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	ent, err := duh.repository.user.findOneByID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = duh.firewall(req, act, ent); err != nil {
		return nil, err
	}

	affected, err := duh.repository.user.deleteOneByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "user does not exists")
		}
		return nil, err
	}

	return &charon.DeleteUserResponse{
		Affected: affected,
	}, nil
}

func (duh *deleteUserHandler) firewall(req *charon.DeleteUserRequest, act *actor, ent *userEntity) error {
	if act.user.id == ent.id {
		return grpc.Errorf(codes.PermissionDenied, "user is not permited to remove himself")
	}
	if act.user.isSuperuser {
		return nil
	}
	if ent.isSuperuser {
		return grpc.Errorf(codes.PermissionDenied, "only superuser can remove other superuser")
	}
	if ent.isStaff {
		switch {
		case act.user.id == ent.createdBy.Int64Or(0):
			if !act.permissions.Contains(charon.UserCanDeleteStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "staff user cannot be removed by owner, missing permission")
			}
			return nil
		case !act.permissions.Contains(charon.UserCanDeleteStaffAsStranger):
			return grpc.Errorf(codes.PermissionDenied, "staff user cannot be removed by stranger, missing permission")
		}
		return nil
	}

	if act.user.id == ent.createdBy.Int64Or(0) {
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
