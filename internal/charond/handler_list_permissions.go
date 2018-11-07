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

type listPermissionsHandler struct {
	*handler
}

func (lph *listPermissionsHandler) List(ctx context.Context, req *charonrpc.ListPermissionsRequest) (*charonrpc.ListPermissionsResponse, error) {
	act, err := lph.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = lph.firewall(req, act); err != nil {
		return nil, err
	}

	entities, err := lph.repository.permission.Find(ctx, &model.PermissionFindExpr{
		Offset:  req.Offset.Int64Or(0),
		Limit:   req.Limit.Int64Or(10),
		OrderBy: mapping.OrderBy(req.OrderBy),
		Where: &model.PermissionCriteria{
			Subsystem: req.Subsystem,
			Module:    req.Module,
			Action:    req.Action,
		},
	})
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "find permission query failed", err)
	}

	permissions := make([]string, 0, len(entities))
	for _, e := range entities {
		permissions = append(permissions, e.Permission().String())
	}
	return &charonrpc.ListPermissionsResponse{
		Permissions: permissions,
	}, nil
}

func (lph *listPermissionsHandler) firewall(req *charonrpc.ListPermissionsRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.PermissionCanRetrieve) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "list of permissions cannot be retrieved, missing permission")
}
