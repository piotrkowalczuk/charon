package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type deleteUserHandler struct {
	*handler
}

func (duh *deleteUserHandler) handle(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	duh.loggerWith("user_id", req.Id)

	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "charon: user cannot be deleted, invalid id: %d", req.Id)
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

	affected, err := duh.repository.user.DeleteByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &DeleteUserResponse{
		Affected: affected,
	}, nil
}

func (duh *deleteUserHandler) firewall(req *DeleteUserRequest, act *actor, ent *userEntity) error {
	if act.user.ID == ent.ID {
		return grpc.Errorf(codes.PermissionDenied, "charon: user is not permited to remove himself")
	}
	if act.user.IsSuperuser {
		return nil
	}
	if ent.IsSuperuser {
		return grpc.Errorf(codes.PermissionDenied, "charon: only superuser can remove other superuser")
	}
	if ent.IsStaff {
		switch {
		case act.user.ID == ent.CreatedBy.Int64Or(0):
			if !act.permissions.Contains(UserCanDeleteStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "charon: staff user cannot be removed by owner, missing permission")
			}
			return nil
		case !act.permissions.Contains(UserCanDeleteStaffAsStranger):
			return grpc.Errorf(codes.PermissionDenied, "charon: staff user cannot be removed by stranger, missing permission")
		}
		return nil
	}

	if act.user.ID == ent.CreatedBy.Int64Or(0) {
		if !act.permissions.Contains(UserCanDeleteAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "charon: user cannot be removed by owner, missing permission")
		}
		return nil
	}
	if !act.permissions.Contains(UserCanDeleteAsStranger) {
		return grpc.Errorf(codes.PermissionDenied, "charon: user cannot be removed by stranger, missing permission")
	}
	return nil
}
