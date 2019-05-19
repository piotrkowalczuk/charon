package charond

import (
	"context"

	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"

	"google.golang.org/grpc/codes"
)

type setGroupPermissionsHandler struct {
	*handler
}

func (sgph *setGroupPermissionsHandler) SetPermissions(ctx context.Context, req *charonrpc.SetGroupPermissionsRequest) (*charonrpc.SetGroupPermissionsResponse, error) {
	act, err := sgph.Actor(ctx)
	if err != nil {
		return nil, err
	}

	if err = sgph.firewall(req, act); err != nil {
		return nil, err
	}

	permissions := charon.NewPermissions(req.Permissions...)
	if req.Force {
		_, err := sgph.repository.permission.InsertMissing(ctx, permissions)
		if err != nil {
			return nil, err
		}
	}

	created, removed, err := sgph.repository.group.SetPermissions(ctx, req.GroupId, permissions...)
	if err != nil {
		switch model.ErrorConstraint(err) {
		case model.TableGroupPermissionsConstraintGroupIDForeignKey:
			return nil, grpcerr.E(codes.NotFound, "%s: group does not exist", err.(*pq.Error).Detail)
		case model.TableGroupPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey:
			return nil, grpcerr.E(codes.NotFound, "%s: permission does not exist", err.(*pq.Error).Detail)
		default:
			return nil, err
		}
	}

	return &charonrpc.SetGroupPermissionsResponse{
		Created:   created,
		Removed:   removed,
		Untouched: untouched(int64(len(req.Permissions)), created, removed),
	}, nil
}

func (sgph *setGroupPermissionsHandler) firewall(req *charonrpc.SetGroupPermissionsRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupPermissionCanCreate) && act.Permissions.Contains(charon.GroupPermissionCanDelete) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "group permissions cannot be set, missing permission")
}
