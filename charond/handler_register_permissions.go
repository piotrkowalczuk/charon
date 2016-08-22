package charond

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type registerPermissionsHandler struct {
	*handler
	registry permissionRegistry
}

func (rph *registerPermissionsHandler) handle(ctx context.Context, req *charon.RegisterPermissionsRequest) (*charon.RegisterPermissionsResponse, error) {
	permissions := charon.NewPermissions(req.Permissions...)
	created, untouched, removed, err := rph.registry.register(permissions)
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
		Created:   created,
		Untouched: untouched,
		Removed:   removed,
	}, nil
}
