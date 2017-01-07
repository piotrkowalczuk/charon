package charond

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type deleteGroupHandler struct {
	*handler
}

func (dgh *deleteGroupHandler) Delete(ctx context.Context, req *charonrpc.DeleteGroupRequest) (*wrappers.BoolValue, error) {
	act, err := dgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = dgh.firewall(req, act); err != nil {
		return nil, err
	}

	aff, err := dgh.repository.group.DeleteOneByID(ctx, req.Id)
	if err != nil {
		switch model.ErrorConstraint(err) {
		case model.TableUserGroupsConstraintGroupIDForeignKey:
			return nil, grpc.Errorf(codes.FailedPrecondition, "group cannot be removed, is not empty")
		default:
			return nil, err
		}
	}

	if aff == 0 {
		return nil, grpc.Errorf(codes.NotFound, "group does not exists")
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

	return grpc.Errorf(codes.PermissionDenied, "group cannot be removed, missing permission")
}
