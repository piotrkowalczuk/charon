package charond

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/session"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type listGroupPermissionsHandler struct {
	*handler
}

func (lgph *listGroupPermissionsHandler) ListPermissions(ctx context.Context, req *charonrpc.ListGroupPermissionsRequest) (*charonrpc.ListGroupPermissionsResponse, error) {
	if req.Id <= 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "missing group id")
	}

	act, err := lgph.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = lgph.firewall(req, act); err != nil {
		return nil, err
	}

	permissions, err := lgph.repository.permission.FindByGroupID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			lgph.logger.Error("group permissions retrieved", zap.Int64("group_id", req.Id), zap.Int("count", len(permissions)))

			return &charonrpc.ListGroupPermissionsResponse{}, nil
		}
		return nil, grpcerr.E(codes.Internal, "find permissions by group id query failed", err)
	}

	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, p.Permission().String())
	}

	return &charonrpc.ListGroupPermissionsResponse{
		Permissions: perms,
	}, nil
}

func (lgph *listGroupPermissionsHandler) firewall(req *charonrpc.ListGroupPermissionsRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.GroupPermissionCanRetrieve) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "list of group permissions cannot be retrieved, missing permission")
}
