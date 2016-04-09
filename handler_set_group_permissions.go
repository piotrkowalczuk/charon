package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setGroupPermissionsHandler struct {
	*handler
}

func (sgph *setGroupPermissionsHandler) handle(ctx context.Context, req *SetGroupPermissionsRequest) (*SetGroupPermissionsResponse, error) {
	act, err := sgph.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	if err = sgph.firewall(req, act); err != nil {
		return nil, err
	}

	created, removed, err := sgph.repository.group.SetPermissions(req.GroupId, NewPermissions(req.Permissions...)...)
	if err != nil {
		return nil, err
	}

	return &SetGroupPermissionsResponse{
		Created:   created,
		Removed:   removed,
		Untouched: untouched(int64(len(req.Permissions)), created, removed),
	}, nil
}

func (sgph *setGroupPermissionsHandler) firewall(req *SetGroupPermissionsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(GroupPermissionCanCreate) && act.permissions.Contains(GroupPermissionCanDelete) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charon: group permissions cannot be set, missing permission")
}
