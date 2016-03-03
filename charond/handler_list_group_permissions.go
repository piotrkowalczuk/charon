package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listGroupPermissionsHandler struct {
	*handler
}

func (lgph *listGroupPermissionsHandler) handle(ctx context.Context, req *charon.ListGroupPermissionsRequest) (*charon.ListGroupPermissionsResponse, error) {
	lgph.loggerWith("group_id", req.Id)

	permissions, err := lgph.repository.permission.FindByUserID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			sklog.Debug(lgph.logger, "user permissions retrieved", "user_id", req.Id, "count", len(permissions))

			return &charon.ListUserPermissionsResponse{}, nil
		}
		return nil, err
	}

	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, p.Permission().String())
	}

	lgph.loggerWith("results", len(permissions))

	return &charon.ListUserPermissionsResponse{
		Permissions: perms,
	}, nil
}
