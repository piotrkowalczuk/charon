package charond

import (
	"database/sql"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
)

type listUserPermissionsHandler struct {
	*handler
}

func (luph *listUserPermissionsHandler) ListPermissions(ctx context.Context, req *charonrpc.ListUserPermissionsRequest) (*charonrpc.ListUserPermissionsResponse, error) {
	permissions, err := luph.repository.permission.FindByUserID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			sklog.Debug(luph.logger, "User Permissions retrieved", "user_id", req.Id, "count", len(permissions))

			return &charonrpc.ListUserPermissionsResponse{}, nil
		}
		return nil, err
	}

	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, p.Permission().String())
	}

	return &charonrpc.ListUserPermissionsResponse{
		Permissions: perms,
	}, nil
}

func (luph *listUserPermissionsHandler) firewall(req *charonrpc.ListUserPermissionsRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.UserPermissionCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "list of User Permissions cannot be retrieved, missing permission")
}
