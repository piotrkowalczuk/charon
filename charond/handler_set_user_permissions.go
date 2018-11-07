package charond

import (
	"context"

	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"

	"google.golang.org/grpc/codes"
)

type setUserPermissionsHandler struct {
	*handler
}

func (suph *setUserPermissionsHandler) SetPermissions(ctx context.Context, req *charonrpc.SetUserPermissionsRequest) (*charonrpc.SetUserPermissionsResponse, error) {
	act, err := suph.Actor(ctx)
	if err != nil {
		return nil, err
	}

	if err = suph.firewall(req, act); err != nil {
		return nil, err
	}

	permissions := charon.NewPermissions(req.Permissions...)
	if req.Force {
		_, err := suph.repository.permission.InsertMissing(ctx, permissions)
		if err != nil {
			return nil, err
		}
	}

	created, removed, err := suph.repository.user.SetPermissions(ctx, req.UserId, permissions...)
	if err != nil {
		switch model.ErrorConstraint(err) {
		case model.TableUserPermissionsConstraintUserIDForeignKey:
			return nil, grpcerr.E(codes.NotFound, "%s: user does not exist", err.(*pq.Error).Detail)
		case model.TableUserPermissionsConstraintPermissionSubsystemPermissionModulePermissionActionForeignKey:
			return nil, grpcerr.E(codes.NotFound, "%s: permission does not exist", err.(*pq.Error).Detail)
		default:
			return nil, err
		}
	}

	return &charonrpc.SetUserPermissionsResponse{
		Created:   created,
		Removed:   removed,
		Untouched: untouched(int64(len(req.Permissions)), created, removed),
	}, nil
}

func (suph *setUserPermissionsHandler) firewall(req *charonrpc.SetUserPermissionsRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}

	if act.Permissions.Contains(charon.UserPermissionCanCreate) && act.Permissions.Contains(charon.UserPermissionCanDelete) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "user permissions cannot be set, missing permission")
}
