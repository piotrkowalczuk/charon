package charond

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/sklog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listGroupPermissionsHandler struct {
	*handler
}

func (luph *listGroupPermissionsHandler) ListPermissions(ctx context.Context, req *charonrpc.ListGroupPermissionsRequest) (*charonrpc.ListGroupPermissionsResponse, error) {
	permissions, err := luph.repository.permission.FindByGroupID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			sklog.Debug(luph.logger, "group permissions retrieved", "group_id", req.Id, "count", len(permissions))

			return &charonrpc.ListGroupPermissionsResponse{}, nil
		}
		return nil, err
	}

	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, p.Permission().String())
	}

	return &charonrpc.ListGroupPermissionsResponse{
		Permissions: perms,
	}, nil
}

func (luph *listGroupPermissionsHandler) firewall(req *charonrpc.ListGroupPermissionsRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupPermissionCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "list of group permissions cannot be retrieved, missing permission")
}
