package main

import (
	"database/sql"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
)

// SetUserPermissions implements charon.RPCServer interface.
func (rs *rpcServer) SetUserPermissions(ctx context.Context, req *charon.SetUserPermissionsRequest) (*charon.SetUserPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: set user permissions endpoint is not implemented yet")
}

// ListUserPermissions implements charon.RPCServer interface.
func (rs *rpcServer) ListUserPermissions(ctx context.Context, req *charon.ListUserPermissionsRequest) (*charon.ListUserPermissionsResponse, error) {
	permissions, err := rs.repository.permission.FindByUserID(int64(req.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			sklog.Debug(rs.logger, "user permissions retrieved", "user_id", int64(req.Id), "count", len(permissions))

			return &charon.ListUserPermissionsResponse{}, nil
		}
		return nil, err
	}

	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, p.Permission().String())
	}

	sklog.Debug(rs.logger, "user permissions retrieved", "user_id", int64(req.Id), "count", len(permissions))

	return &charon.ListUserPermissionsResponse{
		Permissions: perms,
	}, nil
}
