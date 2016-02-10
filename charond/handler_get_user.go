package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getUserHandler struct {
	*handler
}

func (guh *getUserHandler) handle(ctx context.Context, req *charon.GetUserRequest) (*charon.GetUserResponse, error) {
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

	return &charon.GetUserResponse{
		User: ent.Message(),
	}, nil
}

func (guh *getUserHandler) firewall(req *charon.GetUserRequest, act *actor, ent *userEntity) error {
	if act.user.IsSuperuser {
		return nil
	}
	if ent.IsSuperuser {
		return grpc.Errorf(codes.PermissionDenied, "charond: only superuser is permited to retrieve other superuser")
	}
	if ent.IsStaff {
		if ent.CreatedBy.Int64Or(0) == act.user.ID {
			if !act.permissions.Contains(charon.UserCanRetrieveStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "charond: staff user cannot be retrieved as an owner, missing permission")
			}
			return nil
		}
		if !act.permissions.Contains(charon.UserCanRetrieveStaffAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "charond: staff user cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if ent.CreatedBy.Int64Or(0) == act.user.ID {
		if !act.permissions.Contains(charon.UserCanRetrieveAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "charond: user cannot be retrieved as an owner, missing permission")
		}
		return nil
	}
	if !act.permissions.Contains(charon.UserCanRetrieveAsStranger) {
		return grpc.Errorf(codes.PermissionDenied, "charond: user cannot be retrieved as a stranger, missing permission")
	}
	return nil
}
