package main

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// GetPermission implements charon.RPCServer interface.
func (rs *rpcServer) GetPermission(ctx context.Context, req *charon.GetPermissionRequest) (*charon.GetPermissionResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: get permission endpoint is not implemented yet")
}

// RegisterPermissions implements charon.RPCServer interface.
func (rs *rpcServer) RegisterPermissions(ctx context.Context, req *charon.RegisterPermissionsRequest) (*charon.RegisterPermissionsResponse, error) {
	permissions := charon.NewPermissions(req.Permissions...)
	created, untouched, removed, err := rs.permissionRegistry.Register(permissions)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	sklog.Debug(rs.logger, "permission registration success",
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

// ListPermissions implements charon.RPCServer interface.
func (rs *rpcServer) ListPermissions(ctx context.Context, req *charon.ListPermissionsRequest) (*charon.ListPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: list permissions endpoint is not implemented yet")
}
