package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

type registerPermissionsHandler struct {
	*handler
	registry model.PermissionRegistry
}

func (rph *registerPermissionsHandler) Register(ctx context.Context, req *charonrpc.RegisterPermissionsRequest) (*charonrpc.RegisterPermissionsResponse, error) {
	permissions := charon.NewPermissions(req.Permissions...)
	created, untouched, removed, err := rph.registry.Register(ctx, permissions)
	if err != nil {
		switch err {
		case model.ErrEmptySliceOfPermissions, model.ErrEmptySubsystem, model.ErrorInconsistentSubsystem:
			return nil, errf(codes.InvalidArgument, err.Error())
		default:
			return nil, err
		}
	}

	return &charonrpc.RegisterPermissionsResponse{
		Created:   created,
		Untouched: untouched,
		Removed:   removed,
	}, nil
}
