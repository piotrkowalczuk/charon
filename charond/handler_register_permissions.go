package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type registerPermissionsHandler struct {
	*handler
	registry model.PermissionRegistry
}

func (rph *registerPermissionsHandler) Register(ctx context.Context, req *charonrpc.RegisterPermissionsRequest) (*charonrpc.RegisterPermissionsResponse, error) {
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

	return &charonrpc.RegisterPermissionsResponse{
		Created:   created,
		Untouched: untouched,
		Removed:   removed,
	}, nil
}
