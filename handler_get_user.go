package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getUserHandler struct {
	*handler
}

func (guh *getUserHandler) handle(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	guh.loggerWith("user_id", req.Id)

	act, err := guh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	ent, err := guh.repository.user.FindOneByID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = guh.firewall(req, act, ent); err != nil {
		return nil, err
	}

	return guh.response(ent)
}

func (guh *getUserHandler) firewall(req *GetUserRequest, act *actor, ent *userEntity) error {
	if act.user.IsSuperuser {
		return nil
	}
	if ent.IsSuperuser {
		return grpc.Errorf(codes.PermissionDenied, "charon: only superuser is permited to retrieve other superuser")
	}
	if ent.IsStaff {
		if ent.CreatedBy.Int64Or(0) == act.user.ID {
			if !act.permissions.Contains(UserCanRetrieveStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "charon: staff user cannot be retrieved as an owner, missing permission")
			}
			return nil
		}
		if !act.permissions.Contains(UserCanRetrieveStaffAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "charon: staff user cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if ent.CreatedBy.Int64Or(0) == act.user.ID {
		if !act.permissions.Contains(UserCanRetrieveAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "charon: user cannot be retrieved as an owner, missing permission")
		}
		return nil
	}
	if !act.permissions.Contains(UserCanRetrieveAsStranger) {
		return grpc.Errorf(codes.PermissionDenied, "charon: user cannot be retrieved as a stranger, missing permission")
	}
	return nil
}

func (guh *getUserHandler) response(ent *userEntity) (*GetUserResponse, error) {
	msg, err := ent.Message()
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{
		User: msg,
	}, nil
}
