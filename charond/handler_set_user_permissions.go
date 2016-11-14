package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setUserPermissionsHandler struct {
	*handler
}

func (suph *setUserPermissionsHandler) SetPermissions(ctx context.Context, req *charonrpc.SetUserPermissionsRequest) (*charonrpc.SetUserPermissionsResponse, error) {
	act, err := suph.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	if err = suph.firewall(req, act); err != nil {
		return nil, err
	}

	created, removed, err := suph.repository.user.SetPermissions(req.UserId, charon.NewPermissions(req.Permissions...)...)
	if err != nil {
		return nil, err
	}

	return &charonrpc.SetUserPermissionsResponse{
		Created:   created,
		Removed:   removed,
		Untouched: untouched(int64(len(req.Permissions)), created, removed),
	}, nil
}

func (suph *setUserPermissionsHandler) firewall(req *charonrpc.SetUserPermissionsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}

	if act.permissions.Contains(charon.UserPermissionCanCreate) && act.permissions.Contains(charon.UserPermissionCanDelete) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "user permissions cannot be set, missing permission")
}
