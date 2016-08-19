package charond

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setGroupPermissionsHandler struct {
	*handler
}

func (sgph *setGroupPermissionsHandler) handle(ctx context.Context, req *charon.SetGroupPermissionsRequest) (*charon.SetGroupPermissionsResponse, error) {
	act, err := sgph.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	if err = sgph.firewall(req, act); err != nil {
		return nil, err
	}

	created, removed, err := sgph.repository.group.setPermissions(req.GroupId, charon.NewPermissions(req.Permissions...)...)
	if err != nil {
		return nil, err
	}

	return &charon.SetGroupPermissionsResponse{
		Created:   created,
		Removed:   removed,
		Untouched: untouched(int64(len(req.Permissions)), created, removed),
	}, nil
}

func (sgph *setGroupPermissionsHandler) firewall(req *charon.SetGroupPermissionsRequest, act *actor) error {
	if act.user.isSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.GroupPermissionCanCreate) && act.permissions.Contains(charon.GroupPermissionCanDelete) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "group permissions cannot be set, missing permission")
}
