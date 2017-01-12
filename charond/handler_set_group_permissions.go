package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
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

	created, removed, err := sgph.repository.group.SetPermissions(ctx, req.GroupId, charon.NewPermissions(req.Permissions...)...)
	if err != nil {
		return nil, err
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

	return grpc.Errorf(codes.PermissionDenied, "group Permissions cannot be set, missing permission")
}
