package charond

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getUserHandler struct {
	*handler
}

func (guh *getUserHandler) Get(ctx context.Context, req *charonrpc.GetUserRequest) (*charonrpc.GetUserResponse, error) {
	act, err := guh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	ent, err := guh.repository.user.FindOneByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "user does not exists")
		}
		return nil, err
	}
	if err = guh.firewall(req, act, ent); err != nil {
		return nil, err
	}

	return guh.response(ent)
}

func (guh *getUserHandler) firewall(req *charonrpc.GetUserRequest, act *actor, ent *model.UserEntity) error {
	if act.user.IsSuperuser {
		return nil
	}
	if ent.IsSuperuser {
		return grpc.Errorf(codes.PermissionDenied, "only superuser is permited to retrieve other superuser")
	}
	if ent.IsStaff {
		if ent.CreatedBy.Int64Or(0) == act.user.ID {
			if !act.permissions.Contains(charon.UserCanRetrieveStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "staff user cannot be retrieved as an owner, missing permission")
			}
			return nil
		}
		if !act.permissions.Contains(charon.UserCanRetrieveStaffAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "staff user cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if ent.CreatedBy.Int64Or(0) == act.user.ID {
		if !act.permissions.Contains(charon.UserCanRetrieveAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "user cannot be retrieved as an owner, missing permission")
		}
		return nil
	}
	if !act.permissions.Contains(charon.UserCanRetrieveAsStranger) {
		return grpc.Errorf(codes.PermissionDenied, "user cannot be retrieved as a stranger, missing permission")
	}
	return nil
}

func (guh *getUserHandler) response(ent *model.UserEntity) (*charonrpc.GetUserResponse, error) {
	msg, err := ent.Message()
	if err != nil {
		return nil, err
	}
	return &charonrpc.GetUserResponse{
		User: msg,
	}, nil
}
