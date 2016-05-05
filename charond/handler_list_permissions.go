package charond

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listPermissionsHandler struct {
	*handler
}

func (lph *listPermissionsHandler) handle(ctx context.Context, req *charon.ListPermissionsRequest) (*charon.ListPermissionsResponse, error) {
	act, err := lph.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = lph.firewall(req, act); err != nil {
		return nil, err
	}

	entities, err := lph.repository.permission.Find(&permissionCriteria{
		offset:    req.Offset.Int64Or(0),
		limit:     req.Limit.Int64Or(10),
		subsystem: req.Subsystem,
		module:    req.Module,
		action:    req.Action,
	})
	if err != nil {
		return nil, err
	}

	permissions := make([]string, 0, len(entities))
	for _, e := range entities {
		permissions = append(permissions, e.Permission().String())
	}
	return &charon.ListPermissionsResponse{
		Permissions: permissions,
	}, nil
}

func (lph *listPermissionsHandler) firewall(req *charon.ListPermissionsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.PermissionCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charond: list of permissions cannot be retrieved, missing permission")
}
