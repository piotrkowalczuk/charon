package charond

import (
	"context"

	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
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
			return nil, grpcerr.E(codes.InvalidArgument, err)
		default:
			return nil, grpcerr.E(codes.Internal, "permission registration failure", err)
		}
	}

	return &charonrpc.RegisterPermissionsResponse{
		Created:   created,
		Untouched: untouched,
		Removed:   removed,
	}, nil
}
