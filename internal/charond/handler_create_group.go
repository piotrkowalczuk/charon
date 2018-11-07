package charond

import (
	"context"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"google.golang.org/grpc/codes"
)

type createGroupHandler struct {
	*handler
}

func (cgh *createGroupHandler) Create(ctx context.Context, req *charonrpc.CreateGroupRequest) (*charonrpc.CreateGroupResponse, error) {
	if len(req.Name) < 3 {
		return nil, grpcerr.E(codes.InvalidArgument, "group name is required and needs to be at least 3 characters long")
	}
	act, err := cgh.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = cgh.firewall(req, act); err != nil {
		return nil, err
	}

	ent, err := cgh.repository.group.Create(ctx, act.User.ID, req.Name, req.Description)
	if err != nil {
		switch model.ErrorConstraint(err) {
		case model.TableGroupConstraintNameUnique:
			return nil, grpcerr.E(codes.AlreadyExists, "group with given name already exists")
		default:
			return nil, grpcerr.E(codes.Internal, "group fetch failure", err)
		}
	}

	return cgh.response(ent)
}

func (cgh *createGroupHandler) firewall(req *charonrpc.CreateGroupRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupCanCreate) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "group cannot be created, missing permission")
}

func (cgh *createGroupHandler) response(ent *model.GroupEntity) (*charonrpc.CreateGroupResponse, error) {
	msg, err := mapping.ReverseGroup(ent)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "group entity mapping failure", err)
	}
	return &charonrpc.CreateGroupResponse{Group: msg}, nil
}
