package charond

import (
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"google.golang.org/grpc/codes"
)

type deleteGroupHandler struct {
	*handler
}

func (dgh *deleteGroupHandler) Delete(ctx context.Context, req *charonrpc.DeleteGroupRequest) (*wrappers.BoolValue, error) {
	if req.Id <= 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "group cannot be deleted, invalid id")
	}

	act, err := dgh.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = dgh.firewall(req, act); err != nil {
		return nil, err
	}

	aff, err := dgh.repository.group.DeleteOneByID(ctx, req.Id)
	if err != nil {
		//if err == sql.ErrNoRows {
		//	return nil, grpcerr.E(codes.NotFound, "group deletion failure", err)
		//}
		switch model.ErrorConstraint(err) {
		case model.TableUserGroupsConstraintGroupIDForeignKey:
			return nil, grpcerr.E(codes.FailedPrecondition, "group cannot be removed, users are assigned to it")
		case model.TableGroupPermissionsConstraintGroupIDForeignKey:
			return nil, grpcerr.E(codes.FailedPrecondition, "group cannot be removed, permissions are assigned to it")
		default:
			return nil, grpcerr.E(codes.Internal, "group deletion failure", err)
		}
	}

	if aff == 0 {
		return nil, grpcerr.E(codes.NotFound, "group cannot be removed, does not exists")
	}

	return &wrappers.BoolValue{
		Value: aff > 0,
	}, nil
}

func (dgh *deleteGroupHandler) firewall(req *charonrpc.DeleteGroupRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupCanDelete) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "group cannot be removed, missing permission")
}
