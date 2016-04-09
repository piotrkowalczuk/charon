package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setUserGroupsHandler struct {
	*handler
}

func (sugh *setUserGroupsHandler) handle(ctx context.Context, req *SetUserGroupsRequest) (*SetUserGroupsResponse, error) {
	sugh.loggerWith("user_id", req.UserId)

	act, err := sugh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	if err = sugh.firewall(req, act); err != nil {
		return nil, err
	}

	created, removed, err := sugh.repository.userGroups.Set(req.UserId, req.Groups)
	if err != nil {
		return nil, err
	}

	return &SetUserGroupsResponse{
		Created:   created,
		Removed:   removed,
		Untouched: untouched(int64(len(req.Groups)), created, removed),
	}, nil
}

func (sugh *setUserGroupsHandler) firewall(req *SetUserGroupsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(UserGroupCanCreate) && act.permissions.Contains(UserGroupCanDelete) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charon: user groups cannot be set, missing permission")
}
