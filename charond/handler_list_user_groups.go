package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func (lugh *listUserGroupsHandler) firewall(req *charon.ListUserGroupsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.UserGroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charond: list of user groups cannot be retrieved, missing permission")
}
