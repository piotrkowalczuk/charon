package main

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
)

type listGroupPermissionsHandler struct {
	*handler
}

func (lgph *listGroupPermissionsHandler) handle(ctx context.Context, req *charon.ListGroupPermissionsRequest) (*charon.ListGroupPermissionsResponse, error) {
	lgph.loggerWith("group_id", req.Id)

	permissions, err := lgph.repository.permission.FindByUserID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			sklog.Debug(lgph.logger, "group permissions retrieved", "user_id", req.Id, "count", len(permissions))

			return &charon.ListGroupPermissionsResponse{}, nil
		}
		return nil, err
	}

	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, p.Permission().String())
	}

	lgph.loggerWith("results", len(permissions))

	return &charon.ListGroupPermissionsResponse{
		Permissions: perms,
	}, nil
}
