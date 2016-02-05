package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type registerPermissionsHandler struct {
	*handler
	registry PermissionRegistry
}

func (rph *registerPermissionsHandler) handle(ctx context.Context, req *charon.RegisterPermissionsRequest) (*charon.RegisterPermissionsResponse, error) {
	permissions := charon.NewPermissions(req.Permissions...)
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

	return &charon.RegisterPermissionsResponse{
		Created:   int32(created),
		Untouched: int32(untouched),
		Removed:   int32(removed),
	}, nil
}
