package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
)

type listUserGroupsHandler struct {
	*handler
}

func (lugh *listUserGroupsHandler) handle(ctx context.Context, req *charon.ListUserGroupsRequest) (*charon.ListUserGroupsResponse, error) {
	lugh.loggerWith("user_id", req.Id)

	entities, err := lugh.repository.group.FindByUserID(req.Id)
	if err != nil {
		return nil, err
	}

	groups := make([]*charon.Group, 0, len(entities))
	for _, e := range entities {
		groups = append(groups, e.Message())
	}

	return &charon.ListUserGroupsResponse{
		Groups: groups,
	}, nil
}
