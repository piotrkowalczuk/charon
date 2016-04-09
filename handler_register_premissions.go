package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type registerPermissionsHandler struct {
	*handler
	registry PermissionRegistry
}

func (rph *registerPermissionsHandler) handle(ctx context.Context, req *RegisterPermissionsRequest) (*RegisterPermissionsResponse, error) {
	permissions := NewPermissions(req.Permissions...)
	created, untouched, removed, err := rph.registry.Register(permissions)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	rph.loggerWith(
		"registrants", permissions[0].Subsystem(),
		"created", created,
		"untouched", untouched,
		"removed", removed,
		"count", len(req.Permissions),
	)

	return &RegisterPermissionsResponse{
		Created:   created,
		Untouched: untouched,
		Removed:   removed,
	}, nil
}
