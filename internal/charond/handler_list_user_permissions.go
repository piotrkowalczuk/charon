package charond

import (
	"context"

	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/session"
	"google.golang.org/grpc/codes"
)

type listUserPermissionsHandler struct {
	*handler
}

func (luph *listUserPermissionsHandler) ListPermissions(ctx context.Context, req *charonrpc.ListUserPermissionsRequest) (*charonrpc.ListUserPermissionsResponse, error) {
	if req.Id <= 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "missing user id")
	}
	act, err := luph.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = luph.firewall(req, act); err != nil {
		return nil, err
	}

	permissions, err := luph.repository.permission.FindByUserID(ctx, req.Id)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "find permissions by user id query failed", err)
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
	if act.User.ID == req.Id {
		return nil
	}
	if act.Permissions.Contains(charon.UserPermissionCanRetrieve) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "list of user permissions cannot be retrieved, missing permission")
}
