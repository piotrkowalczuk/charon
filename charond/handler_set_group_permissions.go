package charond

import (
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setGroupPermissionsHandler struct {
	*handler
}

func (sgph *setGroupPermissionsHandler) SetPermissions(ctx context.Context, req *charonrpc.SetGroupPermissionsRequest) (*charonrpc.SetGroupPermissionsResponse, error) {
	act, err := sgph.retrieveActor(ctx)
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
			return nil, errf(codes.NotFound, "%s: group does not exist", err.(*pq.Error).Detail)
		case model.TableGroupPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey:
			return nil, errf(codes.NotFound, "%s: permission does not exist", err.(*pq.Error).Detail)
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

	return grpc.Errorf(codes.PermissionDenied, "group permissions cannot be set, missing permission")
}
